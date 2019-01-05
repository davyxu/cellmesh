package main

import (
	"github.com/davyxu/cellmesh/discovery/memsd/api"
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
)

var log = golog.New("memsd")

func startSvc() {

	config := memsd.DefaultConfig()
	if *flagAddr != "" {
		config.Address = *flagAddr
	}

	model.Queue = cellnet.NewEventQueue()
	model.Queue.EnableCapturePanic(true)
	model.Queue.StartLoop()

	p := peer.NewGenericPeer("tcp.Acceptor", "memsd", config.Address, model.Queue)

	msgFunc := proto.GetMessageHandler("memsd")

	proc.BindProcessorHandler(p, "memsd.svc", func(ev cellnet.Event) {

		if msgFunc != nil {
			msgFunc(ev)
		}
	})

	p.(cellnet.TCPSocketOption).SetSocketBuffer(1024*1024, 1024*1024, true)
	p.(cellnet.PeerCaptureIOPanic).EnableCaptureIOPanic(true)
	p.Start()
	service.WaitExitSignal()
}

func broadcast(peer cellnet.Peer, msg interface{}) {
	peer.(cellnet.TCPAcceptor).VisitSession(func(ses cellnet.Session) bool {
		ses.Send(msg)
		return true
	})
}

func deleteValue(ses cellnet.Session, key, reason string) {
	valueMeta := model.DeleteValue(key)

	var ack proto.ValueDeleteNotifyACK
	ack.Key = key

	if valueMeta != nil {
		ack.SvcName = valueMeta.SvcName
	}

	if valueMeta != nil {

		if valueMeta.SvcName == "" {
			log.Infof("DeleteValue '%s'  reason: %s", key, reason)
		} else {
			log.Infof("DeregisterService '%s'  reason: %s", key, reason)
		}
	}

	broadcast(ses.Peer(), &ack)

	ses.Send(&proto.DeleteValueACK{
		Key: key,
	})
}

func init() {

	proto.Handle_Memsd_SetValueREQ = func(ev cellnet.Event) {
		msg := ev.Message().(*proto.SetValueREQ)

		meta := &model.ValueMeta{
			Key:   msg.Key,
			Value: msg.Value,
		}

		// 注册服务
		if msg.SvcName != "" {
			meta.SvcName = msg.SvcName
			meta.Ses = ev.Session()
		}

		model.SetValue(msg.Key, meta)

		if msg.SvcName == "" {
			log.Infof("SetValue '%s' value(size:%d)", msg.Key, len(msg.Value))
		} else {
			log.Infof("RegisterService '%s'", meta.ValueAsServiceDesc().ID)
		}

		broadcast(ev.Session().Peer(), &proto.ValueChangeNotifyACK{
			Key:     msg.Key,
			Value:   msg.Value,
			SvcName: msg.SvcName,
		})

		ev.Session().Send(&proto.SetValueACK{})

	}

	proto.Handle_Memsd_GetValueREQ = func(ev cellnet.Event) {
		msg := ev.Message().(*proto.GetValueREQ)

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

		deleteValue(ev.Session(), msg.Key, "api")
	}

	proto.Handle_Memsd_PullValueREQ = func(ev cellnet.Event) {

		model.VisitValue(func(meta *model.ValueMeta) bool {

			ev.Session().Send(&proto.ValueChangeNotifyACK{
				Key:     meta.Key,
				Value:   meta.Value,
				SvcName: meta.SvcName,
			})

			return true

		})

		ev.Session().Send(&proto.PullValueACK{})
	}

	proto.Handle_Memsd_ClearSvcREQ = func(ev cellnet.Event) {

		log.Infoln("ClearSvc")

		var svcToDelete []*model.ValueMeta
		model.VisitValue(func(meta *model.ValueMeta) bool {

			if meta.SvcName != "" {
				svcToDelete = append(svcToDelete, meta)
			}

			return true
		})

		for _, meta := range svcToDelete {
			deleteValue(meta.Ses, meta.Key, "clearsvc")
		}

		ev.Session().Send(&proto.ClearSvcACK{})
	}

	proto.Handle_Memsd_ClearKeyREQ = func(ev cellnet.Event) {

		log.Infoln("ClearKey")

		var svcToDelete []*model.ValueMeta
		model.VisitValue(func(meta *model.ValueMeta) bool {

			if meta.SvcName == "" {
				svcToDelete = append(svcToDelete, meta)
			}

			return true
		})

		for _, meta := range svcToDelete {
			deleteValue(meta.Ses, meta.Key, "clearkey")
		}

		ev.Session().Send(&proto.ClearKeyACK{})
	}

	proto.Handle_Memsd_Default = func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionClosed:

			var svcToDelete []*model.ValueMeta
			model.VisitValue(func(meta *model.ValueMeta) bool {

				if meta.Ses == ev.Session() {

					// 工具写入的db服务，要持久化保存

					if meta.ValueAsServiceDesc().GetMeta("@Persist") == "" {
						svcToDelete = append(svcToDelete, meta)
					}
				}

				return true
			})

			for _, meta := range svcToDelete {
				deleteValue(meta.Ses, meta.Key, "offline")
			}

		}
	}
}
