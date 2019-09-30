package main

import (
	"github.com/davyxu/cellmesh/svc/memsd/model"
	"github.com/davyxu/cellmesh/svc/memsd/proto"
	"github.com/davyxu/cellnet"
	"strconv"
)

func msgHandler(ev cellnet.Event) {
	switch ev.Message().(type) {
	case *sdproto.AuthREQ:
		auth(ev)
	case *sdproto.ClearKeyREQ:
		clearSvc(ev)
	case *sdproto.ClearSvcREQ:
		clearKey(ev)
	case *sdproto.DeleteValueREQ:
		deleteValue(ev)
	case *sdproto.GetValueREQ:
		getValue(ev)
	case *sdproto.SetValueREQ:
		setValue(ev)
	case *cellnet.SessionClosed:
		sessionClose(ev)
	}
}

func setValue(ev cellnet.Event) {
	msg := ev.Message().(*sdproto.SetValueREQ)

	if !checkAuth(ev.Session()) {

		ev.Session().Send(&sdproto.SetValueACK{
			Code: sdproto.ResultCode_Result_AuthRequire,
		})
		return
	}

	meta := &model.ValueMeta{
		Key:   msg.Key,
		Value: msg.Value,
	}

	// 注册服务
	if model.IsServiceKey(msg.Key) {
		meta.SvcName = msg.SvcName
		meta.Token = model.GetSessionToken(ev.Session())
	}

	model.SetValue(msg.Key, meta)

	if model.IsServiceKey(msg.Key) {
		log.Infof("RegisterService '%s'", meta.ValueAsServiceDesc().ID)
	} else {
		log.Infof("SetValue '%s' value(size:%d)", msg.Key, len(msg.Value))
	}

	model.Broadcast(&sdproto.ValueChangeNotifyACK{
		Key:     msg.Key,
		Value:   msg.Value,
		SvcName: msg.SvcName,
	})

	ev.Session().Send(&sdproto.SetValueACK{})
}

func getValue(ev cellnet.Event) {
	msg := ev.Message().(*sdproto.GetValueREQ)

	if !checkAuth(ev.Session()) {

		ev.Session().Send(&sdproto.GetValueACK{
			Code: sdproto.ResultCode_Result_AuthRequire,
		})
		return
	}

	valueMeta := model.GetValue(msg.Key)
	if valueMeta != nil {
		ev.Session().Send(&sdproto.GetValueACK{
			Key:   msg.Key,
			Value: valueMeta.Value,
		})
	} else {
		ev.Session().Send(&sdproto.GetValueACK{
			Key:  msg.Key,
			Code: sdproto.ResultCode_Result_NotExists,
		})
	}
}

func deleteValue(ev cellnet.Event) {
	msg := ev.Message().(*sdproto.DeleteValueREQ)

	if !checkAuth(ev.Session()) {

		ev.Session().Send(&sdproto.DeleteValueACK{
			Code: sdproto.ResultCode_Result_AuthRequire,
		})
		return
	}

	deleteValueRecurse(msg.Key, "api")

	ev.Session().Send(&sdproto.DeleteValueACK{
		Key: msg.Key,
	})
}

func auth(ev cellnet.Event) {
	msg := ev.Message().(*sdproto.AuthREQ)

	// 下发所有的值做缓冲
	model.VisitValue(func(meta *model.ValueMeta) bool {

		ev.Session().Send(&sdproto.ValueChangeNotifyACK{
			Key:     meta.Key,
			Value:   meta.Value,
			SvcName: meta.SvcName,
		})

		return true

	})

	var ack sdproto.AuthACK

	// 首次生成token并与ses绑定
	if msg.Token == "" {
		ack.Token = strconv.Itoa(int(model.IDGen.Generate()))
	}

	ev.Session().(cellnet.ContextSet).SetContext("token", ack.Token)

	ev.Session().Send(&ack)
}

func clearSvc(ev cellnet.Event) {
	if !checkAuth(ev.Session()) {
		ev.Session().Send(&sdproto.ClearSvcACK{
			Code: sdproto.ResultCode_Result_AuthRequire,
		})
		return
	}

	log.Infoln("ClearSvc")

	var svcToDelete []*model.ValueMeta
	model.VisitValue(func(meta *model.ValueMeta) bool {

		if meta.SvcName != "" {
			svcToDelete = append(svcToDelete, meta)
		}

		return true
	})

	for _, meta := range svcToDelete {
		deleteNotify(meta.Key, "clearsvc")
	}

	ev.Session().Send(&sdproto.ClearSvcACK{})
}

func clearKey(ev cellnet.Event) {
	if !checkAuth(ev.Session()) {
		ev.Session().Send(&sdproto.ClearKeyACK{
			Code: sdproto.ResultCode_Result_AuthRequire,
		})
		return
	}

	log.Infoln("ClearValue")

	var svcToDelete []*model.ValueMeta
	model.VisitValue(func(meta *model.ValueMeta) bool {

		if meta.SvcName == "" {
			svcToDelete = append(svcToDelete, meta)
		}

		return true
	})

	for _, meta := range svcToDelete {
		deleteNotify(meta.Key, "clearkey")
	}

	ev.Session().Send(&sdproto.ClearKeyACK{})
}

func sessionClose(ev cellnet.Event) {
	if checkAuth(ev.Session()) {
		var svcToDelete []*model.ValueMeta
		model.VisitValue(func(meta *model.ValueMeta) bool {

			if meta.Token == model.GetSessionToken(ev.Session()) {

				// 工具写入的db服务，要持久化保存

				if meta.ValueAsServiceDesc().GetMeta("@Persist") == "" {
					svcToDelete = append(svcToDelete, meta)
				}
			}

			return true
		})

		for _, meta := range svcToDelete {
			deleteNotify(meta.Key, "offline")
		}
	}
}
