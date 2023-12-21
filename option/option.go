package main

const metaFileName = "plume_log_202306091128.log"

//func main() {
//	//var (
//	//	err          error
//	//	openMetaFile *os.File
//	//)
//
//	removeLogfFile := fmt.Sprintf("plume_log_%s.log", time.Now().Add(-1*time.Minute).Format("200601021504"))
//	open, _ := os.Open(removeLogfFile)
//	//defer func() {
//	closeErr := open.Close()
//	if closeErr != nil {
//		log.Error(closeErr.Error())
//	}
//	log.Info("close file: %s", removeLogfFile)
//	//}()
//	defer func() {
//		removeErr := os.Remove(removeLogfFile)
//		if removeErr != nil {
//			log.Error(removeErr.Error())
//		}
//	}()
//}
