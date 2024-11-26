package util

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/pbgen/protos_animal"
	"github.com/yaoguangduan/proto-editor/pbgen/protos_user"
	"github.com/yaoguangduan/proto-editor/pbgen/protos_world"
	"google.golang.org/protobuf/encoding/protojson"
	"testing"
)

func TestMessageList(t *testing.T) {
	user := protos_user.NewUserEditor()
	user.AddFriends(protos_user.NewFriendEditor().SetName("name"))
	user.AddFriends(protos_user.NewFriendEditor().SetName("name3"))
	user.GetFriends().Get(1).SetRelation(32.1)
	user.AddTempChange(23).AddTempChange(33).AddTempChange(43)

	ori := &protos_user.User{}
	user.CopyTo(ori)
	another := protos_user.NewUserEditor()
	another.CopyFrom(ori)
	another.FlushDirty(false)
	user.FlushDirty(false)

	user.GetFriends().Get(0).SetRelation(333)
	user.GetFriends().Get(1).SetRelation(66)
	user.GetFriends().Get(1).SetName("mabaa")
	user.GetFriends().Add(protos_user.NewFriendEditor().SetName("name2"))
	user.GetFriends().Get(2).SetName("new name").SetRelation(999)
	user.GetFriends().RemoveByIdx(0)

	user.GetTempChange().Add(53)
	user.GetTempChange().Remove(23)
	user.GetTempChange().Set(0, 83)
	user.GetTempChange().Clear()

	fmt.Println(user.GetFriends().ChangeLog())
	fmt.Println(user.GetFriends().MonitorDirtyList())

	dto := &protos_user.User{}
	user.GatherDirtyToPb(dto)
	fmt.Println(protojson.Format(dto))

	another.MergeDirtyFromPb(dto)
	fmt.Println(another.GetFriends().ChangeLog())
	fmt.Println(another.GetFriends().MonitorDirtyList())
	ur := &protos_user.User{}
	user.CopyTo(ur)
	ar := &protos_user.User{}
	another.CopyTo(ar)
	fmt.Println(protojson.Format(ar))
	fmt.Println(protojson.Format(ur))

}

func TestBasicMap(t *testing.T) {
	world := protos_world.NewWorldEditor()
	world.GetTreeCntMap().Put("bai", 2)
	world.GetTreeCntMap().Put("song", 10)
	world.GetAnimals().Put(999, protos_animal.NewAnimalEditor().SetType(protos_animal.AnimalType_AT_Cat).SetCount(23))
	world.GetAnimals().Put(888, protos_animal.NewAnimalEditor().SetType(protos_animal.AnimalType_AT_Elephant).SetCount(43))
	another := protos_world.NewWorldEditor()
	original := &protos_world.World{}
	world.CopyTo(original)
	another.CopyFrom(original)
	world.FlushDirty(false)
	another.FlushDirty(false)

	world.SetId(999999)
	world.GetTreeCntMap().Put("liu", 22)
	world.GetTreeCntMap().Remove("bai")
	world.GetTreeCntMap().Remove("ssss")
	world.GetAnimals().Put(777, protos_animal.NewAnimalEditor().SetType(protos_animal.AnimalType_AT_Dog).SetCount(94))
	world.GetAnimals().Get(999).SetCount(89)
	world.GetAnimals().Remove(888)

	dirty := &protos_world.World{}
	world.GatherDirtyToPb(dirty)
	another.MergeDirtyFromPb(dirty)

	wo := &protos_world.World{}
	world.CopyTo(wo)
	ao := &protos_world.World{}
	another.CopyTo(ao)
	fmt.Println(protojson.MarshalOptions{Indent: ""}.Format(wo))
	fmt.Println(protojson.MarshalOptions{Indent: ""}.Format(ao))
}
