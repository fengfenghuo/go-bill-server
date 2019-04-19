package cluster

import (
	"database/sql"

	"github.com/bill-server/go-bill-server/node/common"
	"github.com/bill-server/go-bill-server/node/sfoxdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cluster")

type SFoxCluster struct {
	db     *sql.DB
	pool   *redis.Pool
	height int64
	stmts  map[string]*sql.Stmt
}

func NewSFoxCluster(dbUri string) (*SFoxCluster, error) {
	db, err := sql.Open("mysql", dbUri)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.SetConnMaxLifetime(60000)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	cluster := SFoxCluster{
		db:    db,
		stmts: map[string]*sql.Stmt{},
	}
	// if err = cluster.initDb(); err != nil {
	// 	return nil, err
	// }

	// if err = cluster.initStmts(); err != nil {
	// 	return nil, err
	// }
	// if err = cluster.initHeight(); err != nil {
	// 	return nil, err
	// }
	//初始化redis连接池
	/*pool := cluster.newPool("127.0.0.1:6379", "")
	conn := pool.Get()
	_, err = conn.Do("SET", "pool", "test")
	defer conn.Close()
	if err != nil {
		log.Error("redis:redis initialization faile,err" + err.Error())
		return nil, err
	}
	address := common.AccountID{}
	err = address.UnmarshalText([]byte("0x06a7e7203e869ad407c8a637fcb8c9e15f28f09e"))
	if err != nil {

	}
	account := sfoxdb.Account{Address: address}
	account.Asset.SFox = 10000
	account.Asset.NU = 10000000
	account.LastTxOut.Index = 12
	//conn.
	// &account.Asset.SFox, &account.Asset.NU, &block_id, &account.LastTxOut.Index)
	err = cluster.setAccount(account)
	if err != nil {

	}
	cluster.getAccount("0x06a7e7203e869ad407c8a637fcb8c9e15f28f09e")*/
	return &cluster, nil
}

func (cluster *SFoxCluster) CurentHeight() (height int64, err error) {
	return cluster.height, nil
}

func (cluster *SFoxCluster) QueryTXByID(id int64) (tx *sfoxdb.Tx, err error) {
	return nil, nil
}

func (cluster *SFoxCluster) QueryAccountByAddress(address common.AccountID) (account *sfoxdb.Account, err error) {
	return nil, nil
}
