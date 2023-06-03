package store

import (
	"context"
	"encoding/json"
	"fmt"
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
func (ps *Store) GetGroup(id string) (*Group, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(constructKey2(id), nil)
	if err != nil {
		return nil, err
	}

	post := &Group{}
	for _, pair := range data {
		post := Group{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}

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
func (ps *Store) GetAllGroups(ctx context.Context) ([]*Group, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(allGroups, nil)
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

	return posts, nil
}

func (ps *Store) Delete(id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructKey(id, version, ""), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}
func (ps *Store) DeleteGroup(ctx context.Context, id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.DeleteTree(constructKey2(id), nil)
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

	sid, rid := generateKey(post.Version, post.Labels)
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

	data, _, err := kv.List(constructKey(id, version, labels), nil)
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
