package cluster

var stmtDefins = map[string]string{
	"stmt_account_insert":      "  INSERT INTO account (account_id)VALUES(?) ",
	"stmt_account_update":      "  UPDATE account SET up_time = NOW() WHERE account_id=? ",
	"stmt_account_query":       "  SELECT * FROM account where account_id=? ",
	"stmt_tx_insert":           "  INSERT INTO tx (tx_id,account_id)VALUES(?,?) ",
	"stmt_tx_query_by_id":      "  SELECT * FROM tx where tx_id=? ",
	"stmt_tx_query_by_account": "  SELECT * FROM tx where account_id=? ",
}

func (cluster *SFoxCluster) initStmts() error {
	for name, def := range stmtDefins {
		stmt, err := cluster.db.Prepare(def)
		if err != nil {
			log.Error("Cluster: prepare stmt " + name + " error, " + err.Error() + ", [" + def + "]")
			return err
		}
		cluster.stmts[name] = stmt
	}
	return nil
}
