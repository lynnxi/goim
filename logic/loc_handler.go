package main

// import (
// 	. "game-im/lib/mio"
// )

// func GpackLocHandler(req *Gpack) (res *Gpack, err error) {
// 	from := req.Gameid
// 	lat := req.Body.(map[string]interface{})["Lat"].(float64)
// 	lng := req.Body.(map[string]interface{})["Lng"].(float64)
// 	ext := req.Body.(map[string]interface{})["Ext"].(string)

// 	//gid := req.Body.(map[string]interface{})["gid"].(string)
// 	gid := "1"
// 	sendLpsh(from, gid, req.Appid, lat, lng, ext)

// 	res = NewRetGpack(req.Gameid, req.Appid, req.Sid, map[string]string{"ec": "0"})

// 	return
// }

// func sendLpsh(from string, gid string, appid string, lat float64, lng float64, ext string) {
// 	gameids := getGameidsByGid(gid, appid)

// 	var locs []interface{}
// 	loc := map[string]interface{}{
// 		"Fr":  from,
// 		"Lat": lat,
// 		"Lng": lng,
// 		"Ext": ext,
// 	}

// 	locs = append(locs, loc)

// 	for _, gameid := range gameids {
// 		if gameid != from { //不给自己发
// 			gpack := NewReqGpack(gameid, appid, "loc_psh", map[string]interface{}{"gid": gid, "locs": locs})
// 			logic.OutputPackBufferChannel <- gpack
// 		}
// 	}
// }
