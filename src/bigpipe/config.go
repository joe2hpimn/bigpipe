package bigpipe

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type TopicInfo struct {
	Partitions int
}

type ProducerACL struct {
	Secret string
	Topic string
	Name string
}

type Config struct {
	Kafka_bootstrap_servers string

	Kafka_producer_channel_size int
	Kafka_producer_retries int

	// topic信息
	Kafka_producer_topics map[string]TopicInfo
	// acl访问权限
	Kafka_producer_acl map[string]ProducerACL

	Http_server_port int
	Http_server_read_timeout int
	Http_server_write_timeout int
}

var config Config

func LoadConfig(path string) bool {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	dict := map[string]interface{} {}

	err = json.Unmarshal(content, &dict)
	if err != nil {
		return false
	}

	config.Kafka_bootstrap_servers = dict["kafka.bootstrap.servers"].(string)
	config.Kafka_producer_channel_size = int(dict["kafka.producer.channel.size"].(float64))
	config.Kafka_producer_retries = int(dict["kafka.producer.retries"].(float64))

	config.Http_server_port = int(dict["http.server.port"].(float64))
	config.Http_server_read_timeout = int(dict["http.server.read.timeout"].(float64))
	config.Http_server_write_timeout = int(dict["http.server.write.timeout"].(float64))

	config.Kafka_producer_topics = map[string]TopicInfo{}

	topicsArr := dict["kafka.producer.topics"].([]interface{})
	for _, value := range topicsArr {
		topicMap := value.(map[string]interface{})
		name := topicMap["name"].(string)
		partitions := int(topicMap["partitions"].(float64))
		config.Kafka_producer_topics[name] = TopicInfo{Partitions: partitions}
	}

	config.Kafka_producer_acl = map[string]ProducerACL{}

	aclArr := dict["kafka.producer.acl"].([]interface{})
	for _, value := range aclArr {
		aclMap := value.(map[string]interface{})
		name := aclMap["name"].(string)
		secret := aclMap["secret"].(string)
		topic := aclMap["topic"].(string)
		config.Kafka_producer_acl[name] = ProducerACL{Name: name, Secret: secret, Topic: topic}
		// 检查acl涉及的topic是否配置
		if _, exists := config.Kafka_producer_topics[topic]; !exists {
			fmt.Println("ACL中配置的topic: " + topic + " 不存在,请检查kafka.producer.topics.")
			return false
		}
	}

	return true
}

func GetConfig() *Config {
	return &config
}