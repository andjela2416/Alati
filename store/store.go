package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"os"
)

type Store struct {
	cli *api.Client
}

func New() (*Store, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Store{
		cli: client,
	}, nil
}

func (ps *Store) Get(id string, version string) ([]*Config, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	configs := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
func (ps *Store) GetGroup(id string, version string) ([]*Group, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructGroupKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	configs := []*Group{}
	for _, pair := range data {
		config := &Group{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
func (ps *Store) GetOneGroup(id string, version string) (*Group, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructGroupKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	for _, pair := range data {
		config := &Group{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	// Ako nijedna grupa nije pronađena, možete vratiti odgovarajuću grešku
	return nil, errors.New("group not found, not exist")
}
func (ps *Store) GetOneConfig(id string, version string) (*Config, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	// Ako nijedna grupa nije pronađena, možete vratiti odgovarajuću grešku
	return nil, errors.New("config not found, not exist")
}

func (cs *Store) SaveGroup(post *Group) (*Group, error) {
	kv := cs.cli.KV()

	sid, rid := generateGroupKey(post.Version, post.Labels)
	post.Id = rid

	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data} // constructKey2(post.Id), Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *Store) GetAll() ([]*Config, error) {

	kv := ps.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {

		return nil, err
	}

	configs := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {

			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}
func (ps *Store) GetAllGroups() ([]*Group, error) {

	kv := ps.cli.KV()
	data, _, err := kv.List(allGroups, nil)
	if err != nil {

		return nil, err
	}

	groups := []*Group{}
	for _, pair := range data {
		group := &Group{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {

			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (ps *Store) Delete(id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}
func (ps *Store) DeleteGroup(id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructGroupKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *Store) Config(config *Config) (*Config, error) {

	kv := ps.cli.KV()

	sid, rid := generateKey(config.Version, config.Labels)
	config.Id = rid

	data, err := json.Marshal(config)
	if err != nil {

		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {

		return nil, err
	}

	return config, nil
}

func (ps *Store) PostGroup(post *Group) (*Group, error) {

	kv := ps.cli.KV()

	sid, rid := generateGroupKey(post.Version, post.Labels)
	post.Id = rid

	data, err := json.Marshal(post)
	if err != nil {

		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {

		return nil, err
	}

	return post, nil
}

func (ps *Store) GetGroupsByLabels(id string, version string, labels string) ([]*Group, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructGroupKey(id, version, labels), nil)
	if err != nil {
		return nil, err
	}

	posts := []*Group{}

	for _, pair := range data {
		post := &Group{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}
func (ps *Store) GetConfigsByLabels(id string, version string, labels string) ([]*Config, error) {
	kv := ps.cli.KV()

	data, _, err := kv.List(constructKey(id, version, labels), nil)
	if err != nil {
		return nil, err
	}

	configs := []*Config{}

	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	if err != nil {
		return nil, err
	}

	return configs, nil
}

func (ps *Store) SaveRequestId() string {
	kv := ps.cli.KV()

	reqId := generateRequestId()

	i := &api.KVPair{Key: reqId, Value: nil}

	_, err := kv.Put(i, nil)
	if err != nil {
		return "error"
	}
	return reqId

}

func (ps *Store) FindRequestId(requestId string) bool {
	kv := ps.cli.KV()

	key, _, err := kv.Get(requestId, nil)

	fmt.Println(key)

	if err != nil || key == nil {
		return false
	}

	return true

}
func generateRequestId() string {
	rid := uuid.New().String()

	return rid
}

/*
func (cs *Store) PutGroupConfigByGroupId(config *Config, groupId string) (*Group, error) {
	kv := cs.cli.KV()

	data,_,err := kv.List(constructGroupKey(groupId,"",""),nil)
	if err != nil {

	}

	var newConfigGroup *Group
	for _,pair := range data{
		configGroup :=&Group{}
		err = json.Unmarshal(pair.Value,configGroup)
		if err!=nil{}
		_,err := kv.DeleteTree(constructGroupKey(groupId,"",""),nil)
		if err != nil{}

		config.Id= uuid.NewString()

		newConfigGroup=configGroup
		newConfigGroup.Configs=configGroup.Configs
		newConfigGroup.Configs=append(newConfigGroup.Configs,config)

		bid := generateKey(newConfigGroup.Version,"")

		data,err:=json.Marshal(newConfigGroup)
		if err!=nil{}
		p:=&api.KVPair{Key:bid,Value:data}
		_,err =kv.Put(p,nil)
		if err !=nil {}

		sid :=generateGroupKey(groupId,newConfigGroup.Version,config.Id,config.Labels)
		configData,configErr:=json.Marshal(config)
		if configErr !=nil{}

		b:= &api.KVPair{Key:sid,Value:configData}
		_,err =kv.Put(b,nil)
		if err!=nil{}
	}*/
//	return newConfigGroup, nil
//}
