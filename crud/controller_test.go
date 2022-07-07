package crud

import (
	_ "modernc.org/sqlite"
)

// Not use

// func test_init() {
// 	db.Info = config.DatabaseInfoSQLite
// 	db.Info.FilePath = "../test.db"
// 	db.Info.TableName = "books"
// 	err := db.SetupDB()
// 	if err != nil {
// 		log.Fatal("SetupDB:", err)
// 	}
// 	err = db.Obj.CreateDB()
// 	if err != nil {
// 		log.Fatal("CreateDB:", err)
// 	}

// 	// err = db.Obj.CreateTable()
// 	// if err != nil {
// 	// 	log.Fatal("CreateTable:", err)
// 	// }
// }

// func TestInsertData(t *testing.T) {
// 	type args struct {
// 		book model.Book
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		count int64
// 	}{
// 		{
// 			name: "INSERT",
// 			args: args{
// 				book: model.Book{
// 					Title:  null.StringFrom("test_title"),
// 					Author: null.StringFrom("test_author"),
// 				},
// 			},
// 			count: int64(1),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			test_init()
// 			defer os.Remove(db.Info.FilePath)
// 			defer db.Con.Close()

// 			db.Info.TableName = "books"
// 			got, _, err := InsertData(tt.args.book)

// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if got != tt.count {
// 				require.Equal(t, tt.count, got)
// 			}
// 		})
// 	}
// }

// func TestSelectData(t *testing.T) {
// 	type args struct {
// 		id int
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		count int64
// 	}{
// 		{
// 			name: "SELECT",
// 			args: args{id: 0},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			test_init()

// 			db.Info.TableName = "books"
// 			_, err := SelectData(tt.args.id)
// 			defer os.Remove(db.Info.FilePath)
// 			defer db.Con.Close()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		})
// 	}
// }

// func TestUpdateData(t *testing.T) {
// 	type args struct {
// 		book model.Book
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		count int64
// 	}{
// 		{
// 			name: "UPDATE",
// 			args: args{
// 				book: model.Book{
// 					Idx:    null.IntFrom(int64(1)),
// 					Title:  null.StringFrom("test_title"),
// 					Author: null.StringFrom("test_author"),
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			test_init()

// 			db.Info.TableName = "books"
// 			got, err := UpdateData(tt.args.book)
// 			defer os.Remove(db.Info.FilePath)
// 			defer db.Con.Close()
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			if got != tt.count {
// 				require.Equal(t, tt.count, got)
// 			}
// 		})
// 	}
// }

// func TestDeleteData(t *testing.T) {
// 	type args struct {
// 		id int
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		count int64
// 	}{
// 		{
// 			name: "DELETE",
// 			args: args{id: 1},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			test_init()

// 			db.Info.TableName = "books"
// 			_, err := DeleteData(tt.args.id)
// 			defer os.Remove(db.Info.FilePath)
// 			defer db.Con.Close()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		})
// 	}
// }
