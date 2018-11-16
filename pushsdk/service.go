package pushsdk

type PushService interface {
	/**
	 * @title 推送信息标题
	 * @body  推送信息内容
	 * @TRP Token、PushId、RegId的简称
	 */
	DoPush(title string, body string, TPR string) error
}
