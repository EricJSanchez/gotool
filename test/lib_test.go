package test

import (
	"context"
	"fmt"
	"github.com/EricJSanchez/gotool/pkg/environment"
	"github.com/EricJSanchez/gotool/service"
	"github.com/EricJSanchez/gotool/sys"
	"testing"
	"time"
)

func init() {
	// 初始化环境
	sys.InitEnv(environment.Development)
	service.Register()
	//初始化 config
	sys.InitConfig("../configs/")
	//nacos初始化
	_ = service.Factory.Nacos.InitClient(sys.Cfg("app").GetStringMap("nacos"))
	//nacos 注册服务
	//defer service.Factory.Nacos.DeRegister()
	//service.Factory.Nacos.Register()

}

type ChatMessage struct {
	MsgId             string    `gorm:"column:msg_id;not null pk VARCHAR(64)" json:"msg_id"`
	CorpId            string    `gorm:"column:corp_id;not null pk VARCHAR(64)" json:"corp_id"`
	Action            string    `gorm:"column:action;not null VARCHAR(20)" json:"action"`
	From              string    `gorm:"column:from;not null index VARCHAR(100)" json:"from"`
	FromStaffIds      string    `gorm:"column:from_staff_ids;not null index VARCHAR(100)" json:"from_staff_ids"`
	ToList            string    `gorm:"column:to_list;JSONB" json:"to_list"`
	ToListStaffIds    string    `gorm:"column:to_list_staff_ids;JSONB" json:"to_list_staff_ids"`
	RoomId            string    `gorm:"column:room_id;index VARCHAR(64)" json:"room_id"`
	MsgTime           int64     `gorm:"column:msg_time;not null index BIGINT" json:"msg_time"`
	MsgType           string    `gorm:"column:msg_type;not null index VARCHAR(20)" json:"msg_type"`
	Text              string    `gorm:"column:text;JSONB" json:"text"`
	Image             string    `gorm:"column:image;JSONB" json:"image"`
	Revoke            string    `gorm:"column:revoke;JSONB" json:"revoke"`
	Agree             string    `gorm:"column:agree;JSONB" json:"agree"`
	Voice             string    `gorm:"column:voice;JSONB" json:"voice"`
	Video             string    `gorm:"column:video;JSONB" json:"video"`
	Card              string    `gorm:"column:card;JSONB" json:"card"`
	Location          string    `gorm:"column:location;JSONB" json:"location"`
	Emotion           string    `gorm:"column:emotion;JSONB" json:"emotion"`
	File              string    `gorm:"column:file;JSONB" json:"file"`
	Link              string    `gorm:"column:link;JSONB" json:"link"`
	Weapp             string    `gorm:"column:weapp;JSONB" json:"weapp"`
	Chatrecord        string    `gorm:"column:chatrecord;JSONB" json:"chatrecord"`
	Todo              string    `gorm:"column:todo;JSONB" json:"todo"`
	Vote              string    `gorm:"column:vote;JSONB" json:"vote"`
	Redpacket         string    `gorm:"column:redpacket;JSONB" json:"redpacket"`
	Meeting           string    `gorm:"column:meeting;JSONB" json:"meeting"`
	Docmsg            string    `gorm:"column:docmsg;JSONB" json:"docmsg"`
	Markdown          string    `gorm:"column:markdown;JSONB" json:"markdown"`
	News              string    `gorm:"column:news;JSONB" json:"news"`
	CreateTime        time.Time `gorm:"column:create_time;not null DATETIME" json:"create_time"`
	Md5Sums           string    `gorm:"column:md5_sums;VARCHAR(1024)" json:"md5_sums"`
	Calendar          string    `gorm:"column:calendar;JSONB" json:"calendar"`
	Collect           string    `gorm:"column:collect;JSONB" json:"collect"`
	Mixed             string    `gorm:"column:mixed;JSONB" json:"mixed"`
	Voiceid           string    `gorm:"column:voiceid;VARCHAR(200)" json:"voiceid"`
	MeetingVoiceCall  string    `gorm:"column:meeting_voice_call;JSONB" json:"meeting_voice_call"`
	VoipDocShare      string    `gorm:"column:voip_doc_share;JSONB" json:"voip_doc_share"`
	ExternalRedpacket string    `gorm:"column:external_redpacket;JSONB" json:"external_redpacket"`
	SphFeed           string    `gorm:"column:sph_feed;JSONB" json:"sph_feed"`
	//头像和昵称
	Avatar      string `json:"avatar"`
	NickName    string `json:"nick_name"`
	HideMsgTime int    `json:"hide_msg_time"`
	RISKID      string `gorm:"column:risk_id;JSONB" json:"risk_id"`
}

func TestLib(t *testing.T) {
	nacosAddr := sys.Cfg("app").GetString("nacos.addr")
	sys.Pr(nacosAddr)

	sys.Pr(sys.Nacos().GetString("SsoUrl"))

	var str string
	sys.Gorm().Table("ww_staff").Where("id > ?", 1).Limit(1).Select("userid").Find(&str)
	sys.Pr(str)

	rs, _ := sys.Redis().Get(context.Background(), "cms:draw:cache:***").Result()
	sys.Pr(rs)

	result, err := sys.Elastic().Search().
		Index("chat-message-***").
		Pretty(true).
		From(0).
		Size(2).
		Sort("msg_time", false).
		Do(context.Background())
	if err == nil {
		total := int(result.TotalHits())
		res, _ := sys.EsToStruct[ChatMessage](result)
		sys.Pr(res, total)
	} else {
		fmt.Println("es err", err)
	}

	return
}
