package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
)

type News struct {
	ID      int64  `json:"id"`
	Topic   string `json:"topic"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
	Tags    []Tag  `json:"tags"`
}

type Newsfilter struct {
	Status string `json:"status,omitempty"`
	Topic  string `json:"topic,omitempty"`
}

func (n *News) GetNews(db *sql.DB) error {
	var strtag string
	err := db.QueryRow(`
	SELECT news.topic, news.title, news.content, news.status, JSON_ARRAYAGG(JSON_OBJECT('id', tags.id, 'name', tags.name))
	FROM news
	LEFT JOIN (news_tags,tags)ON(news.id = news_tags.news_id AND news_tags.tags_id = tags.id AND tags.status != 'deleted')
	WHERE
		news.status != 'deleted'
		AND news.id = ?
	GROUP BY
		news.id`,
		n.ID).Scan(&n.Topic, &n.Title, &n.Content, &n.Status, &strtag)
	if err != nil {
		return err
	}
	var tmpTag []Tag
	json.Unmarshal([]byte(strtag), &tmpTag)
	if len(tmpTag) > 1 || tmpTag[0].ID != 0 {
		n.Tags = tmpTag
	} else {
		n.Tags = []Tag{}
	}
	return err
}

func (n *News) UpdateNews(db *sql.DB) error {

	BuildUpdateNews := func(news News) string {
		value := reflect.ValueOf(news)
		name := value.Type()
		var newsstr string
		for i := 0; i < value.NumField(); i++ {
			if value.Field(i).Interface() == "" || name.Field(i).Name == "Tags" || name.Field(i).Name == "ID" {
				continue
			} else if newsstr == "" {
				newsstr = fmt.Sprintf("SET %s = '%s'", name.Field(i).Name, value.Field(i).Interface())
			} else {
				newsstr += fmt.Sprintf(", %s = '%s'", name.Field(i).Name, value.Field(i).Interface())
			}
		}
		return newsstr
	}
	sqlQueery := fmt.Sprintf(`UPDATE news %s WHERE id=?`, BuildUpdateNews(*n))
	_, err :=
		db.Exec(sqlQueery, n.ID)
	if err != nil {
		return err
	}
	if len(n.Tags) > 0 {

		db.Exec("DELETE FROM news_tags WHERE news_id=?", n.ID)

		for _, tag := range n.Tags {
			fmt.Println(tag)
			_, err := db.Exec(
				"INSERT INTO news_tags(news_id,tags_id) VALUES(?,(SELECT id from tags WHERE name = ?)) ", n.ID, tag.Name)
			if err != nil {
				return err
			}
		}
	}
	return n.GetNews(db)
}

func (n *News) DeleteNews(db *sql.DB) error {
	_, err := db.Exec("UPDATE news SET status='deleted' WHERE id=?", n.ID)

	return err
}

func (n *News) CreateNews(db *sql.DB) error {

	BuildNews := func(news News) string {
		value := reflect.ValueOf(news)
		name := value.Type()
		var newsstr string
		for i := 0; i < value.NumField(); i++ {
			if value.Field(i).Interface() == "" || name.Field(i).Name == "Tags" || name.Field(i).Name == "ID" {
				continue
			} else if newsstr == "" {
				newsstr = fmt.Sprintf("SET %s = '%s'", name.Field(i).Name, value.Field(i).Interface())
			} else {
				newsstr += fmt.Sprintf(", %s = '%s'", name.Field(i).Name, value.Field(i).Interface())
			}
		}
		return newsstr
	}
	sqlQueery := fmt.Sprintf(`INSERT INTO news %s ;`, BuildNews(*n))
	res, err := db.Exec(sqlQueery)
	if err != nil {
		return err
	}

	n.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	for _, tag := range n.Tags {
		_, err := db.Exec(
			"INSERT INTO news_tags(news_id,tags_id) VALUES(?,(SELECT id from tags WHERE name = ?)) ", n.ID, tag.Name)
		if err != nil {
			n.deleteData(db)
			return err
		}
	}
	return n.GetNews(db)
}

func ListNews(db *sql.DB, filter Newsfilter) ([]News, error) {
	buildFilter := func(filter Newsfilter) string {
		value := reflect.ValueOf(filter)
		name := value.Type()
		var filterstr string
		for i := 0; i < value.NumField(); i++ {
			if value.Field(i).Interface() == "" {
				continue
			}
			if filterstr == "" {
				filterstr = fmt.Sprintf("WHERE news.%s = '%s'", name.Field(i).Name, value.Field(i).Interface())
			} else {
				filterstr += fmt.Sprintf(" AND news.%s = '%s'", name.Field(i).Name, value.Field(i).Interface())
			}
		}
		return filterstr
	}
	sqlQueery := fmt.Sprintf(`
	SELECT news.id, news.topic, news.title, news.content, news.status, JSON_ARRAYAGG(JSON_OBJECT('id', tags.id, 'name', tags.name))
	FROM news
	LEFT JOIN (news_tags,tags)ON(news.id = news_tags.news_id AND news_tags.tags_id = tags.id AND tags.status != 'deleted')
	%s
	GROUP BY
		news.id
	`, buildFilter(filter))
	rows, err := db.Query(sqlQueery)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	news := []News{}

	for rows.Next() {
		var n News
		var strtag string
		if err := rows.Scan(&n.ID, &n.Topic, &n.Title, &n.Content, &n.Status, &strtag); err != nil {
			return nil, err
		}
		var tmpTag []Tag
		json.Unmarshal([]byte(strtag), &tmpTag)
		if len(tmpTag) > 1 || tmpTag[0].ID != 0 {
			n.Tags = tmpTag
		} else {
			n.Tags = []Tag{}
		}

		news = append(news, n)
	}

	return news, nil
}

func (n *News) deleteData(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM news WHERE id=?", n.ID)

	return err
}
