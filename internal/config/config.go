// получение данных из .env или ENV
package config

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

var ErrConfigDataInvalid = errors.New("invalid file data")

// необходимые данные для запуска
type Config struct {
	ServerHost string
	ServerPort string
}

// парсим файл, заполняем поля и проверяем на корректность
func NewConfig(patToFile string) (*Config, error) {
	cfg := &Config{}

	if err := cfg.parse(patToFile); err != nil {
		return nil, err
	}

	cfg.unmarshal()

	if !cfg.valid() {
		return nil, ErrConfigDataInvalid
	}

	return cfg, nil
}

// .env не найден —> идем читать из ENV
// .env найден —> создаем ENV
func (cfg *Config) parse(patToFile string) error {
	if _, err := os.Stat(patToFile); os.IsNotExist(err) {
		log.Print("config: file not found, used ENV")
		return nil
	}

	file, err := os.Open(patToFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)

	for scan.Scan() {
		str := strings.TrimSpace(scan.Text())
		if str == "" {
			continue
		}

		parts := strings.SplitN(str, "=", 2)
		if len(parts) != 2 {
			return ErrConfigDataInvalid
		}

		key, val := parts[0], parts[1]

		if err := os.Setenv(key, val); err != nil {
			return err
		}

		log.Printf("config: setenv - {key}:{value} {%s}:{%s};", key, val)
	}

	log.Print("config: end parse file file")

	return nil
}

func (cfg *Config) unmarshal() {
	cfg.ServerHost = os.Getenv("SERVER_HOST")
	cfg.ServerPort = os.Getenv("SERVER_PORT")
}

// проверяем корректность данных
func (cfg *Config) valid() bool {
	if cfg.ServerHost == "" {
		return false
	}

	port, err := strconv.ParseUint(cfg.ServerPort, 10, 16)
	if err != nil || port == 0 {
		return false
	}

	return true
}
