package main

import (
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
	"github.com/davyxu/cellnet"
	"strconv"
)

func init() {

	proto.Handle_Memsd_SetValueREQ = func(ev cellnet.Event) {
		msg := ev.Message().(*proto.SetValueREQ)

		if !checkAuth(ev.Session()) {

			ev.Session().Send(&proto.SetValueACK{
				Code: proto.ResultCode_Result_AuthRequire,
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

		model.Broadcast(&proto.ValueChangeNotifyACK{
			Key:     msg.Key,
			Value:   msg.Value,
			SvcName: msg.SvcName,
		})

		ev.Session().Send(&proto.SetValueACK{})

	}

	proto.Handle_Memsd_GetValueREQ = func(ev cellnet.Event) {
		msg := ev.Message().(*proto.GetValueREQ)

		if !checkAuth(ev.Session()) {

			ev.Session().Send(&proto.GetValueACK{
				Code: proto.ResultCode_Result_AuthRequire,
			})
			return
		}

		valueMeta := model.GetValue(msg.Key)
		if valueMeta != nil {
			ev.Session().Send(&proto.GetValueACK{
				Key:   msg.Key,
				Value: valueMeta.Value,
			})
		} else {
			ev.Session().Send(&proto.GetValueACK{
				Key:  msg.Key,
				Code: proto.ResultCode_Result_NotExists,
			})
		}

	}

	proto.Handle_Memsd_DeleteValueREQ = func(ev cellnet.Event) {
		msg := ev.Message().(*proto.DeleteValueREQ)

		if !checkAuth(ev.Session()) {

			ev.Session().Send(&proto.DeleteValueACK{
				Code: proto.ResultCode_Result_AuthRequire,
			})
			return
		}

		deleteNotify(msg.Key, "api")

		ev.Session().Send(&proto.DeleteValueACK{
			Key: msg.Key,
		})
	}

	proto.Handle_Memsd_AuthREQ = func(ev cellnet.Event) {

		msg := ev.Message().(*proto.AuthREQ)

		model.VisitValue(func(meta *model.ValueMeta) bool {

			ev.Session().Send(&proto.ValueChangeNotifyACK{
				Key:     meta.Key,
				Value:   meta.Value,
				SvcName: meta.SvcName,
			})

			return true

		})

		var ack proto.AuthACK

		// 首次生成token并与ses绑定
		if msg.Token == "" {
			ack.Token = strconv.Itoa(int(model.IDGen.Generate()))
		}

		ev.Session().(cellnet.ContextSet).SetContext("token", ack.Token)

		ev.Session().Send(&ack)
	}

	proto.Handle_Memsd_ClearSvcREQ = func(ev cellnet.Event) {

		if !checkAuth(ev.Session()) {
			ev.Session().Send(&proto.ClearSvcACK{
				Code: proto.ResultCode_Result_AuthRequire,
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

		ev.Session().Send(&proto.ClearSvcACK{})
	}

	proto.Handle_Memsd_ClearKeyREQ = func(ev cellnet.Event) {

		if !checkAuth(ev.Session()) {
			ev.Session().Send(&proto.ClearKeyACK{
				Code: proto.ResultCode_Result_AuthRequire,
			})
			return
		}

		log.Infoln("ClearKey")

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

		ev.Session().Send(&proto.ClearKeyACK{})
	}

	proto.Handle_Memsd_Default = func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionAccepted:

		case *cellnet.SessionClosed:

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
	}
}
