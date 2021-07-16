package contacts_callback

import (
	"context"
	//	"crypto/aes"
	"crypto/rand"
	"github.com/Peanuttown/dd_contacts/dd_crypto"
	"github.com/Peanuttown/tzzGoUtil/reflect"
	"github.com/Peanuttown/tzzGoUtil/log"
	"github.com/Peanuttown/dd_api"
	"encoding/base64"
	//	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
	//"dingding_data/internal/conf"
)

type Config struct {
	Addr   string `json:"addr"`
	CorpId string `json:"corpId"`
	//AESKey string `json:"aeskey"`
}

// server handle dingding callback
type Server struct {
	server        *http.ServeMux
	logger        log.LoggerLite
	eventHandlers map[dd_api.EventType]EventHandler
	conf          *Config
	ddClient      *dd_api.Client
	aesKey        []byte
	ddCryptIns    *dd_crypto.DingTalkCrypto
	callbackUrl   string
	msgsToHandle chan *msgToHandle
}

func NewServer(
	conf *Config,
	logger log.LoggerLite,
	ddClient *dd_api.Client,
	callbackUrl string,
) *Server {
	if len(conf.CorpId) == 0 {
		panic(fmt.Errorf("CoprId is nil"))
	}
	aesKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		panic(err)
	}
	token := time.Now().String()
	ddCrypt := dd_crypto.NewDingTalkCrypto(token, base64.RawStdEncoding.EncodeToString(aesKey), conf.CorpId)

	retServer := &Server{
		ddCryptIns:    ddCrypt,
		conf:          conf,
		logger:        logger,
		eventHandlers: make(map[dd_api.EventType]EventHandler),
		ddClient:      ddClient,
		aesKey:        aesKey,
		callbackUrl:   callbackUrl,
	}
	retServer.server = http.NewServeMux()
	retServer.server.Handle("/", retServer)

	return retServer
}

func (this *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	logger := this.logger
	logger.Debug("< Get callback msg>")
	etType, reqMap, err := this.resolveEventType(req)
	if err != nil {
		logger.Error(err)
		this.resErr(res)
		return
	}
	logger.Debugf("<Raw params> %+v\n", reqMap)

	if etType == dd_api.EVENT_TYPE_CHECK_URL {
		this.resOk(res)
		return
	}
	// try push to msgs channel and return ok
	select {
	case this.msgsToHandle <- &msgToHandle{EventType:etType,ReqMap:reqMap}:
		logger.Debug("Push Callback msg to task quque successfully")
		this.resOk(res)
	default:
		logger.Error("Task queue has full, respnse error to dingding")
		this.resErr(res)
	}
	return

}

func (this *Server) resOk(res http.ResponseWriter) {
	logger := this.logger
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := make([]byte, 6)
	_, err := rand.Read(nonce)
	if err != nil {
		logger.Error(err)
		return
	}
	nonceStr := base64.StdEncoding.EncodeToString(nonce)
	encryptMsg, sig, err := this.ddCryptIns.GetEncryptMsg("success", timeStamp, nonceStr)
	if err != nil {
		logger.Error(err)
		return
	}
	resEntity := &ResOkEntity{
		Signature:  sig,
		TimeStamp:  timeStamp,
		Nonce:      nonceStr,
		EncryptMsg: encryptMsg,
	}
	resBytes, err := json.Marshal(resEntity)
	if err != nil {
		logger.Error(err)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "application/json")
	res.Write(resBytes)
}

type ResOkEntity struct {
	Signature  string `json:"msg_signature"`
	TimeStamp  string `json:"timeStamp"`
	Nonce      string `json:"nonce"`
	EncryptMsg string `json:"encrypt"`
}

func (this *Server) handleCommonEvent(ctx context.Context, eventType dd_api.EventType, reqMap map[string]interface{}) error {
	logger := this.logger
	handler, ok := this.eventHandlers[eventType]
	if !ok {
		err := fmt.Errorf("接收到未注册的回调事件: %v", eventType)
		logger.Error(err)
		return err
	}
	err := handler(ctx, reqMap)
	if err != nil {
		return err
	}
	return nil
}

func (this *Server) resErr(res http.ResponseWriter) {
	res.WriteHeader(http.StatusInternalServerError)
}

type ReqEntity struct {
	Encrypt string `json:"encrypt"`
}

func (this *Server) resolveEventType(req *http.Request) (etType dd_api.EventType, reqMap map[string]interface{}, err error) {
	logger := this.logger
	// < parse param from url
	signature := req.URL.Query().Get("signature")
	timestamp := req.URL.Query().Get("timestamp")
	nonce := req.URL.Query().Get("nonce")
	// >
	reqBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf("读取请求数据失败: %v", err)
		return "", nil, err
	}
	reqEntity := &ReqEntity{}
	err = json.Unmarshal(reqBytes, reqEntity)
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}
	plainMsg, err := this.ddCryptIns.GetDecryptMsg(signature, timestamp, nonce, reqEntity.Encrypt)
	if err != nil {
		err = fmt.Errorf("< Decrypt Msg failed: <%w>>", err)
		logger.Error(err)
		return "", nil, err
	}
	reqMap = make(map[string]interface{})
	err = json.Unmarshal([]byte(plainMsg), &reqMap)
	if err != nil {
		logger.Errorf("解析请求参数失败: %v", err)
		return "", nil, err
	}
	eventTypeIfce, ok := reqMap["EventType"]
	if !ok {
		err = fmt.Errorf("请求参数中没有 EventType 字段: %+v", reqMap)
		logger.Error(err)
		return "", nil, err
	}

	eventTypeStr, ok := eventTypeIfce.(string)
	if !ok {
		err = fmt.Errorf("eventType 类型是 %s 不是 string类型", reflect.GetStructName(eventTypeIfce))
		logger.Error(err)
		return "", nil, err
	}

	eventType := dd_api.EventType(eventTypeStr)
	return eventType, reqMap, nil
}

func (this *Server) Run(ctx context.Context) error {
	logger := this.logger
	// < 订阅回调事件
	if len(this.eventHandlers) == 0 {
		return fmt.Errorf("回掉函数为空")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	addr := this.conf.Addr
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer l.Close()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				// 判断服务是否启动了
				conn, err := net.Dial("tcp", this.conf.Addr)
				if err != nil {
					logger.Error("服务尚未启动，不发起订阅请求", err)
					continue
				}
				conn.Close()
				// 开始订阅事件
				err = this.subscribeEvent(ctx)
				if err != nil {
					logger.Errorf("订阅事件失败:%v, 退出服务", err)
					cancel()
					return
				}
				//订阅成功
				logger.Debug("订阅事件成功")
				return
			}
		}
	}()
	// < init channel 
	this.msgsToHandle = make(chan *msgToHandle,100)
	// >
	go func(ctx context.Context){ // handle event sequentially
		for{
			select{
			case <-ctx.Done():
				return
			case msg,ok:=<-this.msgsToHandle:
				if !ok {
					continue
				}
				etType := msg.EventType
				reqMap := msg.ReqMap
				logger.Debugf("开始处理钉钉回调事件: %s \n ", etType)
				err = this.handleCommonEvent(ctx, etType, reqMap)
				if err != nil{
					logger.Error("handle msg %+v failed: %w",msg,err)
				}
				// TODO handle msg
			}
		}
	}(ctx)
	go func(){
		<-ctx.Done()
		l.Close()
	}()
	return http.Serve(l, this)
}

func (this *Server) subscribeEvent(ctx context.Context) error {
	etTypes := make([]dd_api.EventType, 0, len(this.eventHandlers))
	for etType, _ := range this.eventHandlers {
		etTypes = append(etTypes, etType)
	}

	err := dd_api.NewApiContactsCallbackRegister(
		this.ddCryptIns.Token,
		etTypes,
		this.ddCryptIns.EncodingAESKey,
		this.callbackUrl,
	).ExecBy(ctx, this.ddClient)
	if err != nil {
		if dd_api.ErrIsCallbackUrlExist(err) {
			// try update
			err = dd_api.NewApiReqCallbackRegUpdate(
				this.ddCryptIns.Token,
				etTypes,
				this.ddCryptIns.EncodingAESKey,
				this.callbackUrl,
			).ExecBy(ctx, this.ddClient)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (this *Server) RegEventUserAddHandler(handle func(ctx context.Context, req *EventUserAddData) error) {
	this.eventHandlers[dd_api.EVENT_TYPE_USER_ADD] = func(ctx context.Context, req map[string]interface{}) error {
		var reqObj = &EventUserAddData{}
		err := unMarshalFromMap(req, reqObj)
		if err != nil {
			return err
		}
		return handle(ctx, reqObj)
	}
}

func (this *Server) RegEventUserLeaveHandler(handle func(ctx context.Context, req *EventUserLeaveData) error) {
	this.eventHandlers[dd_api.EVENT_TYPE_USER_LEAVE] = func(ctx context.Context, req map[string]interface{}) error {
		var reqObj = &EventUserLeaveData{}
		err := unMarshalFromMap(req, reqObj)
		if err != nil {
			return err
		}
		return handle(ctx, reqObj)
	}
}

func (this *Server) RegEventDeptAdd(
	handle func(ctx context.Context, req *EventDeptAdd) error,
) {
	this.eventHandlers[dd_api.EVENT_TYPE_DEPT_ADD] = func(ctx context.Context, req map[string]interface{}) error {
		var reqObj = &EventDeptAdd{}
		err := unMarshalFromMap(req, reqObj)
		if err != nil {
			return err
		}
		return handle(ctx, reqObj)
	}
}

func (this *Server) RegEventDeptDel(
	handle func(ctx context.Context, req *EventDeptDelete) error,
) {
	this.eventHandlers[dd_api.EVENT_TYPE_DEPT_DEL] = func(ctx context.Context, req map[string]interface{}) error {
		var reqObj = &EventDeptDelete{}
		err := unMarshalFromMap(req, reqObj)
		if err != nil {
			return err
		}
		return handle(ctx, reqObj)
	}
}

func (this *Server) RegEventDeptUpdate(
	handle func(ctx context.Context, req *EventDeptUpdate) error,
) {
	this.eventHandlers[dd_api.EVENT_TYPE_DEPT_UPDATE] = func(ctx context.Context, req map[string]interface{}) error {
		var reqObj = &EventDeptUpdate{}
		err := unMarshalFromMap(req, reqObj)
		if err != nil {
			return err
		}
		return handle(ctx, reqObj)
	}
}

func (this *Server) RegEventUserUpdate(
	handle func(ctx context.Context, req *EventUserUpdate) error,
) {
	this.eventHandlers[dd_api.EVENT_TYPE_USER_UPDATE] = func(ctx context.Context, req map[string]interface{}) error {
		var reqObj = &EventUserUpdate{}
		err := unMarshalFromMap(req, reqObj)
		if err != nil {
			return err
		}
		return handle(ctx, reqObj)
	}
}

type msgToHandle struct{
	EventType dd_api.EventType
	ReqMap map[string]interface{}
}

