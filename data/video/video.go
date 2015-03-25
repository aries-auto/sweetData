package video

import (
	"github.com/curt-labs/sweetData/data/brand"
	"github.com/curt-labs/sweetData/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"time"
)

type Video struct {
	ID           int       `json:"id,omitempty" xml:"id,omitempty"`
	Title        string    `json:"title, omitempty" xml:"title,omitempty"`
	VideoType    VideoType `json:"videoType,omitempty" xml:"v,omitempty"`
	Description  string    `json:"description,omitempty" xml:"description,omitempty"`
	DateAdded    time.Time `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	DateModified time.Time `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
	IsPrimary    bool      `json:"isPrimary,omitempty" xml:"v,omitempty"`
	Thumbnail    string    `json:"thumbnail,omitempty" xml:"thumbnail,omitempty"`
	Channels     Channels  `json:"channels,omitempty" xml:"channels,omitempty"`
	Files        CdnFiles  `json:"files,omitempty" xml:"files,omitempty"`
	CategoryIds  []int     `json:"categoryIds,omitempty" xml:"categoryIds,omitempty"`
	PartIds      []int     `json:"partIds,omitempty" xml:"partIds,omitempty"`

	WebsiteId int           `json:"websiteId,omitempty" xml:"websiteId,omitempty"`
	Brands    []brand.Brand `json:"brands,omitempty" xml:"brands,omitempty"`
}
type Videos []Video

type Channel struct {
	ID           int         `json:"id,omitempty" xml:"id,omitempty"`
	Type         ChannelType `json:"type,omitempty" xml:"type,omitempty"`
	Link         string      `json:"link,omitempty" xml:"link,omitempty"`
	EmbedCode    string      `json:"embedCode,omitempty" xml:"embedCode,omitempty"`
	ForiegnID    string      `json:"foreignId,omitempty" xml:"foreignId,omitempty"`
	DateAdded    time.Time   `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	DateModified time.Time   `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
	Title        string      `json:"title,omitempty" xml:"title,omitempty"`
	Description  string      `json:"description,omitempty" xml:"description,omitempty"`
}

type Channels []Channel

type ChannelType struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	Name        string `json:"name,omitempty" xml:"name,omitempty"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
}

type CdnFile struct {
	ID           int         `json:"id,omitempty" xml:"id,omitempty"`
	Type         CdnFileType `json:"type,omitempty" xml:"type,omitempty"`
	Path         string      `json:"path,omitempty" xml:"path,omitempty"`
	Bucket       string      `json:"bucket,omitempty" xml:"bucket,omitempty"`
	ObjectName   string      `json:"objectName,omitempty" xml:"objectName,omitempty"`
	FileSize     string      `json:"fileSize,omitempty" xml:"fileSize,omitempty"`
	DateAdded    time.Time   `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	DateModified time.Time   `json:"dateModified,omitempty" xml:"dateModified,omitempty"`
	LastUploaded string      `json:"lastUploaded,omitempty" xml:"lastUploaded,omitempty"`
}

type CdnFiles []CdnFile

type CdnFileType struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	MimeType    string `json:"mimeType,omitempty" xml:"mimeType,omitempty"`
	Title       string `json:"title,omitempty" xml:"title,omitempty"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
}

type VideoType struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	Icon string `json:"icon,omitempty" xml:"icon,omitempty"`
}

var (
	checkVideo = `select ID from VideoNew where subjectTypeID = ? and title = ? and description = ? and dateAdded = ? and dateModified = ?
		and isPrimary = ? and thumbnail = ? and isPrivate = ?`
	insertVideo = `insert into VideoNew (subjectTypeID, title, description, dateAdded, dateModified, isPrimary, thumbnail, isPrivate)
			values(?,?,?,?,?,?,?,?)`
	checkVideoCdn   = `select ID from VideoCdnFiles where cdnID = ? and videoID = ?`
	checkVideoChan  = `select ID from VideoChannels where videoID = ? and channelID = ?`
	insertVideoCdn  = `insert into VideoCdnFiles (cdnID, videoID) values (?,?)`
	insertVideoChan = `insert into VideoChannels (videoID, channelID) values (?,?)`
)

func (v *Video) Check() (int, error) {
	var err error
	var id int
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkVideo)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		v.VideoType.ID,
		v.Title,
		v.Description,
		v.DateAdded,
		v.DateModified,
		v.IsPrimary,
		v.Thumbnail,
		v.IsPrimary,
	).Scan(&id)
	return id, err
}

func InsertVideos(videos []Video) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertVideo)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, v := range videos {
		id, err := v.Check()
		if id > 0 {
			continue
		}
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		res, err := stmt.Exec(
			v.VideoType.ID,
			v.Title,
			v.Description,
			v.DateAdded,
			v.DateModified,
			v.IsPrimary,
			v.Thumbnail,
			v.IsPrimary,
		)
		if err != nil {
			return err
		}
		vid, err := res.LastInsertId()
		if err != nil {
			return err
		}
		v.ID = int(vid)
		//TODO MAYBE - do we need to check CDN, Chan, VideoType?

		//check + insert CDNs
		for _, cdn := range v.Files {
			joinID, err := cdn.Check(v)
			if joinID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = cdn.Insert(v)
			if err != nil {
				return err
			}
		}

		//check + insert Channels
		for _, channel := range v.Channels {
			joinID, err := channel.Check(v)
			if joinID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = channel.Insert(v)
			if err != nil {
				return err
			}
		}

		//check + insert Brands
		for _, b := range v.Brands {
			joinID, err := b.CheckVideoBrand(v.ID)
			if joinID > 0 {
				continue
			}
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			err = b.InsertVideoBrand(v.ID)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (cdn *CdnFile) Check(v Video) (int, error) {
	var err error
	var id int
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkVideoCdn)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(cdn.ID, v.ID).Scan(&id)
	return id, err
}

func (channel *Channel) Check(v Video) (int, error) {
	var err error
	var id int
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return id, err
	}
	defer db.Close()

	stmt, err := db.Prepare(checkVideoChan)
	if err != nil {
		return id, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(v.ID, channel.ID).Scan(&id)
	return id, err
}

func (cdn *CdnFile) Insert(v Video) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertVideoCdn)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cdn.ID, v.ID)
	return err
}

func (channel *Channel) Insert(v Video) error {
	var err error
	db, err := sql.Open("mysql", database.NewDBConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertVideoChan)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(v.ID, channel.ID)
	return err
}
