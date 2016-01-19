package output

import(
//	"fmt"
//	"errors"
	
	"github.com/gotools/logs"
//	"privatehub"
	"db"
)

func SaveDataToLocalDB(paperTitle string, paperContent string, source_title string, source_address string) {
	
	isExist, err := DataInDB(source_address)
	if err != nil {
		logs.Error("judge url in db or not error")
	} else if isExist == true {
		logs.Debug("current url is exsist in database")
		return;
	}
	
	sql := "insert into zhiliaoyuan(id, paper_title, source_title, source_url, add_time) values("
	sql += "nextval('zhiliaoyuan_id_seq'),'" + paperTitle + "','" 
	sql += source_title + "','" + source_address + "',now());" 

	err = db.ExecuteSql(sql)
	if err != nil {
		logs.Error(err.Error())
	}
}

func DataInDB(url string) (bool, error) {
	sql := "select 1 from zhiliaoyuan where source_url = '" + url +"'";
	return db.DataInDB(sql)
}