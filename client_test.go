package gochatwork

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	testToken := "testToken"

	Convey("normal", t, func() {
		client := New(testToken)
		So(client.config.url, ShouldEqual, chatworkURL)
		So(client.config.token, ShouldEqual, testToken)
	})

	Convey("kddi", t, func() {
		client := NewKddiChatwork(testToken)
		So(client.config.url, ShouldEqual, kddiChatworkURL)
		So(client.config.token, ShouldEqual, testToken)
	})
}

type stubHTTP struct {
	GetCount    int
	GetByte     []byte
	GetEndPoint string
	GetParams   url.Values
}

func (h *stubHTTP) Get(endPoint string, params url.Values, config *config) ([]byte, error) {
	h.GetCount++
	h.GetEndPoint = endPoint
	h.GetParams = params
	return h.GetByte, nil
}

func TestMe(t *testing.T) {
	Convey("correct", t, func() {
		correctJSON := `{
		"account_id":42,
		"room_id":4242,
		"name":"name_text",
		"chatwork_id":"chatwork_id_text",
		"organization_id":424242,
		"organization_name":"organization_name_text",
		"department":"department_text",
		"title":"title_text",
		"url":"url_text",
		"introduction":"introduction_text",
		"mail":"mail_text",
		"tel_organization":"tel_organization_text",
		"tel_extension":"tel_extension_text",
		"tel_mobile":"tel_mobile_text",
		"skype":"skype_text",
		"facebook":"facebook_text",
		"twitter":"twitter_text",
		"avatar_image_url":"avatar_image_url_text"
		}`

		testToken := "testToken"
		client := New(testToken)

		Convey("MeRaw", func() {
			stub := &stubHTTP{}
			stub.GetByte = make([]byte, 0)
			client.connection = stub

			b, _ := client.MeRaw()
			So(len(b), ShouldEqual, 0)
			So(stub.GetCount, ShouldEqual, 1)
			So(stub.GetEndPoint, ShouldEqual, "me")
		})

		Convey("Me", func() {
			stub := &stubHTTP{}
			stub.GetByte = []byte(correctJSON)
			client.connection = stub

			me, err := client.Me()
			So(err, ShouldBeNil)

			So(me.AccountID, ShouldEqual, 42)
			So(me.RoomID, ShouldEqual, 4242)
			So(me.Name, ShouldEqual, "name_text")
			So(me.ChatworkID, ShouldEqual, "chatwork_id_text")
			So(me.OrganizationID, ShouldEqual, 424242)
			So(me.OrganizationName, ShouldEqual, "organization_name_text")
			So(me.Department, ShouldEqual, "department_text")
			So(me.Title, ShouldEqual, "title_text")
			So(me.URL, ShouldEqual, "url_text")
			So(me.Introduction, ShouldEqual, "introduction_text")
			So(me.Mail, ShouldEqual, "mail_text")
			So(me.TelOrganization, ShouldEqual, "tel_organization_text")
			So(me.TelExtension, ShouldEqual, "tel_extension_text")
			So(me.TelMobile, ShouldEqual, "tel_mobile_text")
			So(me.Skype, ShouldEqual, "skype_text")
			So(me.Facebook, ShouldEqual, "facebook_text")
			So(me.Twitter, ShouldEqual, "twitter_text")
			So(me.AvatarImageURL, ShouldEqual, "avatar_image_url_text")
		})
	})

	Convey("connect", t, func() {
		token := os.Getenv("CHATWORK_API_TOKEN")
		if token == "" {
			t.Log("skip this test because no token")
			return
		}

		client := New(token)
		b, err := client.MeRaw()
		So(len(b), ShouldNotEqual, 0)
		So(err, ShouldBeNil)
	})
}

func TestMyStatus(t *testing.T) {
	testToken := "testToken"
	client := New(testToken)

	Convey("correct", t, func() {
		correctJSON := `{
		"unread_room_num":4,
		"mention_room_num":42,
		"mytask_room_num":424,
		"unread_num":4242,
		"mention_num":42424,
		"mytask_num":424242
		}`

		Convey("MyStatus", func() {
			stub := &stubHTTP{}
			stub.GetByte = []byte(correctJSON)
			client.connection = stub

			status, err := client.MyStatus()
			So(err, ShouldBeNil)
			So(status.UnreadRoomNum, ShouldEqual, 4)
			So(status.MentionRoomNum, ShouldEqual, 42)
			So(status.MytaskRoomNum, ShouldEqual, 424)
			So(status.UnreadNum, ShouldEqual, 4242)
			So(status.MentionNum, ShouldEqual, 42424)
			So(status.MytaskNum, ShouldEqual, 424242)
		})

		Convey("MyStatusRaw", func() {
			stub := &stubHTTP{}
			stub.GetByte = make([]byte, 0)
			client.connection = stub

			b, _ := client.MyStatusRaw()
			So(len(b), ShouldEqual, 0)
			So(stub.GetCount, ShouldEqual, 1)
			So(stub.GetEndPoint, ShouldEqual, "my/status")
		})
	})
}

func TestMyTasks(t *testing.T) {
	testToken := "testToken"
	client := New(testToken)

	Convey("correct", t, func() {
		Convey("MyTask", func() {
			correctJSON := `
[
  {
    "task_id":4,
    "room":{
      "room_id":42,
      "name":"room_name_1",
      "icon_path":"room_icon_path_1"
    },
    "assigned_by_account":{
      "account_id":424,
      "name":"assigned_by_account_name_1",
      "avatar_image_url":"assigned_by_account_avatar_image_url_1"
    },
    "message_id":4242,
    "body":"task_body_1",
    "limit_time":42424,
    "status":"done"
  },
  {
    "task_id":424242,
    "room":{
      "room_id":4242424,
      "name":"room_name_2",
      "icon_path":"room_icon_path_2"
    },
    "assigned_by_account":{
      "account_id":42424242,
      "name":"assigned_by_account_name_2",
      "avatar_image_url":"assigned_by_account_avatar_image_url_2"
    },
    "message_id":424242424,
    "body":"task_body_2",
    "limit_time":4242424242,
    "status":"open"
  }
]`
			stub := &stubHTTP{}
			stub.GetByte = []byte(correctJSON)
			client.connection = stub

			params := url.Values{}
			params.Add("assigned_by_account_id", "42")
			params.Add("status", "done")

			tasks, err := client.MyTasks(params)

			So(stub.GetCount, ShouldEqual, 1)
			So(stub.GetEndPoint, ShouldEqual, "my/tasks")
			So(stub.GetParams, ShouldResemble, params)

			So(err, ShouldBeNil)
			So(len(tasks), ShouldEqual, 2)

			task := tasks[0]
			So(task.TaskID, ShouldEqual, 4)
			So(task.Room.RoomID, ShouldEqual, 42)
			So(task.Room.Name, ShouldEqual, "room_name_1")
			So(task.Room.IconPath, ShouldEqual, "room_icon_path_1")
			So(task.AssignedByAccount.AccountID, ShouldEqual, 424)
			So(task.AssignedByAccount.Name, ShouldEqual, "assigned_by_account_name_1")
			So(task.AssignedByAccount.AvatarImageURL, ShouldEqual, "assigned_by_account_avatar_image_url_1")
			So(task.MessageID, ShouldEqual, 4242)
			So(task.Body, ShouldEqual, "task_body_1")
			So(task.LimitTime, ShouldEqual, 42424)
			So(task.Status, ShouldEqual, "done")

			task = tasks[1]
			So(task.TaskID, ShouldEqual, 424242)
			So(task.Room.RoomID, ShouldEqual, 4242424)
			So(task.Room.Name, ShouldEqual, "room_name_2")
			So(task.Room.IconPath, ShouldEqual, "room_icon_path_2")
			So(task.AssignedByAccount.AccountID, ShouldEqual, 42424242)
			So(task.AssignedByAccount.Name, ShouldEqual, "assigned_by_account_name_2")
			So(task.AssignedByAccount.AvatarImageURL, ShouldEqual, "assigned_by_account_avatar_image_url_2")
			So(task.MessageID, ShouldEqual, 424242424)
			So(task.Body, ShouldEqual, "task_body_2")
			So(task.LimitTime, ShouldEqual, 4242424242)
			So(task.Status, ShouldEqual, "open")
		})

		Convey("MyTasksRaw", func() {
			stub := &stubHTTP{}
			stub.GetByte = make([]byte, 0)
			client.connection = stub

			params := url.Values{}
			params.Add("assigned_by_account_id", "42")
			params.Add("status", "done")

			b, _ := client.MyTasksRaw(params)
			So(len(b), ShouldEqual, 0)
			So(stub.GetCount, ShouldEqual, 1)
			So(stub.GetEndPoint, ShouldEqual, "my/tasks")
			So(stub.GetParams, ShouldResemble, params)
		})
	})
}
