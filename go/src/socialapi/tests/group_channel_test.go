package main

import (
	"koding/db/mongodb/modelhelper"
	"math/rand"
	"socialapi/models"
	"socialapi/rest"
	"socialapi/workers/common/runner"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGroupChannel(t *testing.T) {

	r := runner.New("test")
	if err := r.Init(); err != nil {
		t.Fatalf("couldnt start bongo %s", err.Error())
	}
	defer r.Close()

	modelhelper.Initialize(r.Conf.Mongo)
	defer modelhelper.Close()

	Convey("while  testing pinned activity channel", t, func() {
		rand.Seed(time.Now().UnixNano())
		groupName := "testgroup" + strconv.FormatInt(rand.Int63(), 10)
		Convey("channel should be there", func() {

			account := models.NewAccount()
			account.OldId = AccountOldId.Hex()
			account, err := rest.CreateAccount(account)
			So(err, ShouldBeNil)
			So(account, ShouldNotBeNil)

			channel1, err := rest.CreateChannelByGroupNameAndType(account.Id, groupName, models.Channel_TYPE_GROUP)
			So(err, ShouldBeNil)
			So(channel1, ShouldNotBeNil)

			channel2, err := rest.CreateChannelByGroupNameAndType(account.Id, groupName, models.Channel_TYPE_GROUP)
			So(err, ShouldNotBeNil)
			So(channel2, ShouldBeNil)
		})

		Convey("group channel should be shown before announcement", func() {
			account, err := models.CreateAccountInBothDbs()
			So(err, ShouldBeNil)

			channels, err := rest.FetchChannels(account.Id)
			So(err, ShouldBeNil)
			So(len(channels), ShouldEqual, 2)
			So(channels[0].TypeConstant, ShouldEqual, models.Channel_TYPE_GROUP)
			So(channels[1].TypeConstant, ShouldEqual, models.Channel_TYPE_ANNOUNCEMENT)

		})

		Convey("owner should be able to update it", func() {
			account := models.NewAccount()
			account.OldId = AccountOldId.Hex()
			account, err := rest.CreateAccount(account)
			So(err, ShouldBeNil)
			So(account, ShouldNotBeNil)

			channel1, err := rest.CreateChannelByGroupNameAndType(account.Id, groupName, models.Channel_TYPE_GROUP)
			So(err, ShouldBeNil)
			So(channel1, ShouldNotBeNil)
			// fetching channel returns creator id
			_, err = rest.UpdateChannel(channel1)
			So(err, ShouldBeNil)
		})

		Convey("owner should only be able to update name and purpose of the channel", nil)

		Convey("normal user should not be able to update it", func() {
			account := models.NewAccount()
			account.OldId = AccountOldId.Hex()
			account, err := rest.CreateAccount(account)
			So(err, ShouldBeNil)
			So(account, ShouldNotBeNil)

			channel1, err := rest.CreateChannelByGroupNameAndType(account.Id, groupName, models.Channel_TYPE_GROUP)
			So(err, ShouldBeNil)
			So(channel1, ShouldNotBeNil)

			channel1.CreatorId = rand.Int63()
			_, err = rest.UpdateChannel(channel1)
			So(err, ShouldNotBeNil)
		})

		Convey("owner cant delete it", func() {
			account := models.NewAccount()
			account.OldId = AccountOldId.Hex()
			account, err := rest.CreateAccount(account)
			So(err, ShouldBeNil)
			So(account, ShouldNotBeNil)

			channel1, err := rest.CreateChannelByGroupNameAndType(account.Id, groupName, models.Channel_TYPE_GROUP)
			So(err, ShouldBeNil)
			So(channel1, ShouldNotBeNil)

			err = rest.DeleteChannel(account.Id, channel1.Id)
			So(err, ShouldNotBeNil)
		})

		Convey("normal user cant delete it", func() {
			account := models.NewAccount()
			account.OldId = AccountOldId.Hex()
			account, err := rest.CreateAccount(account)
			So(err, ShouldBeNil)
			So(account, ShouldNotBeNil)

			channel1, err := rest.CreateChannelByGroupNameAndType(account.Id, groupName, models.Channel_TYPE_GROUP)
			So(err, ShouldBeNil)
			So(channel1, ShouldNotBeNil)

			err = rest.DeleteChannel(rand.Int63(), channel1.Id)
			So(err, ShouldNotBeNil)
		})

		Convey("member can post status update", nil)

		Convey("non-member can not post status update", nil)
	})
}
