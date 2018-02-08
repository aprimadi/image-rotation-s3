package config

import (
  "flag"
)

type Config struct {
  BaseDir string
  Environment string
}

var Cfg *Config

func ParseConfig() (*Config, []string) {
  baseDir := flag.String("base-dir", ".", "Base directory")
  environmentPtr := flag.String("environment", "development", "Server environment")

  config := &Config{
    BaseDir:      *baseDir,
    Environment:  *environmentPtr,
  }
  args := flag.Args()

  return config, args
}
