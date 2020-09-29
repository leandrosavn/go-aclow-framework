package config

import (
	"os"
	"strconv"
)

func MongoDbDsn() string {
	return os.Getenv("MONGODB_DSN")
}
func MongoDbDatabase() string {
	return os.Getenv("MONGODB_DATABASE")
}
func MongoDbMaxPoolSize() uint64 {
	m := os.Getenv("MONGODB_MAX_POOL_SIZE")
	if m == "" {
		return 100
	}
	n, _ := strconv.Atoi(m)
	return uint64(n)
}

func LDAPBase() string {
	return os.Getenv("LDAP_BASE")
}

func LDAPHost() string {
	return os.Getenv("LDAP_HOST")
}

func LDAPPort() int {
	m := os.Getenv("LDAP_PORT")
	if m == "" {
		return 389
	}
	n, _ := strconv.Atoi(m)
	return n
}

func LDAPUserSSL() bool {
	m := os.Getenv("LDAP_USER_SSL")
	if m == "" {
		return false
	}
	n, _ := strconv.ParseBool(m)
	return n
}

func LDAPBindDN() string {
	return os.Getenv("LDAP_BIND_DN")
}

func LDAPBindPassword() string {
	return os.Getenv("LDAP_BIND_PASSWORD")
}

func LDAPUserFilter() string {
	return os.Getenv("LDAP_USER_FILTER")
}

func LDAPGroupFilter() string {
	return os.Getenv("LDAP_GROUP_FILTER")
}

func LDAPAttributes() string {
	return os.Getenv("LDAP_ATTRIBUTES")
}
