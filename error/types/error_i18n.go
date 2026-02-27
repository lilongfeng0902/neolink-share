package types

import (
	"strings"
)

// ErrorMessage is a map of language code to error message
type ErrorMessage map[string]string

// errorMessages contains localized error messages for each error code
// Supported languages: en (English), zh (Chinese), ja (Japanese), fr (French), ru (Russian), vi (Vietnamese)
var errorMessages = map[ErrorCode]ErrorMessage{
	// General Errors (1xxx)
	ErrorCodeInvalidRequest: {
		"en": "Invalid request parameters",
		"zh": "请求参数无效",
		"ja": "無効なリクエストパラメータ",
		"fr": "Paramètres de requête invalides",
		"ru": "Недействительные параметры запроса",
		"vi": "Tham số yêu cầu không hợp lệ",
	},
	ErrorCodeSensitiveWordsDetected: {
		"en": "Sensitive words detected in content",
		"zh": "内容中检测到敏感词",
		"ja": "コンテンツに敏感な単語が検出されました",
		"fr": "Mots sensibles détectés dans le contenu",
		"ru": "Обнаружены нежелательные слова в контенте",
		"vi": "Phát hiện từ nhạy cảm trong nội dung",
	},
	ErrorCodeViolationFeeGrokCSAM: {
		"en": "Content policy violation detected",
		"zh": "检测到内容违规",
		"ja": "コンテンツポリシー違反が検出されました",
		"fr": "Violation de la politique de contenu détectée",
		"ru": "Обнаружено нарушение политики содержимого",
		"vi": "Phát hiện vi phạm chính sách nội dung",
	},

	// System Errors (2xxx)
	ErrorCodeCountTokenFailed: {
		"en": "Failed to count tokens",
		"zh": "Token 计数失败",
		"ja": "トークン数のカウントに失敗しました",
		"fr": "Échec du comptage des jetons",
		"ru": "Не удалось подсчитать токены",
		"vi": "Không thể đếm token",
	},
	ErrorCodeModelPriceError: {
		"en": "Model pricing configuration error",
		"zh": "模型价格配置错误",
		"ja": "モデル価格設定エラー",
		"fr": "Erreur de configuration des prix du modèle",
		"ru": "Ошибка конфигурации цены модели",
		"vi": "Lỗi cấu hình giá mô hình",
	},
	ErrorCodeInvalidApiType: {
		"en": "Invalid API type",
		"zh": "无效的 API 类型",
		"ja": "無効なAPIタイプ",
		"fr": "Type d'API invalide",
		"ru": "Недействительный тип API",
		"vi": "Loại API không hợp lệ",
	},
	ErrorCodeJsonMarshalFailed: {
		"en": "Failed to marshal JSON",
		"zh": "JSON 序列化失败",
		"ja": "JSONマーシャリングに失敗しました",
		"fr": "Échec du marshaling JSON",
		"ru": "Не удалось упаковать JSON",
		"vi": "Không thể chuyển đổi JSON",
	},
	ErrorCodeJsonUnmarshalFailed: {
		"en": "Failed to unmarshal JSON",
		"zh": "JSON 反序列化失败",
		"ja": "JSONアンマーシャリングに失敗しました",
		"fr": "Échec de l'unmarshaling JSON",
		"ru": "Не удалось распаковать JSON",
		"vi": "Không thể phân tích JSON",
	},
	ErrorCodeDoRequestFailed: {
		"en": "Failed to make HTTP request",
		"zh": "HTTP 请求失败",
		"ja": "HTTPリクエストが失敗しました",
		"fr": "Échec de la requête HTTP",
		"ru": "Не выполнить HTTP-запрос",
		"vi": "Yêu cầu HTTP không thành công",
	},
	ErrorCodeGetChannelFailed: {
		"en": "Failed to get channel information",
		"zh": "获取渠道信息失败",
		"ja": "チャンネル情報の取得に失敗しました",
		"fr": "Échec de la récupération des informations du canal",
		"ru": "Не удалось получить информацию о канале",
		"vi": "Không thể lấy thông tin kênh",
	},
	ErrorCodeGenRelayInfoFailed: {
		"en": "Failed to generate relay information",
		"zh": "生成中继信息失败",
		"ja": "リレー情報の生成に失敗しました",
		"fr": "Échec de la génération des informations de relais",
		"ru": "Не удалось создать информацию о ретрансляции",
		"vi": "Không thể tạo thông tin tiếprelay",
	},

	// Channel Errors (3xxx)
	ErrorCodeChannelNoAvailableKey: {
		"en": "No available API key in channel",
		"zh": "渠道中没有可用的 API 密钥",
		"ja": "チャンネルに利用可能なAPIキーがありません",
		"fr": "Aucune clé API disponible dans le canal",
		"ru": "Нет доступного ключа API в канале",
		"vi": "Không có khóa API khả dụng trong kênh",
	},
	ErrorCodeChannelParamOverrideInvalid: {
		"en": "Invalid channel parameter override",
		"zh": "无效的渠道参数覆盖",
		"ja": "無効なチャンネルパラメータオーバーライド",
		"fr": "Remplacement de paramètre de canal invalide",
		"ru": "Недействительное пер��определение параметра канала",
		"vi": "Ghi đè tham số kênh không hợp lệ",
	},
	ErrorCodeChannelHeaderOverrideInvalid: {
		"en": "Invalid channel header override",
		"zh": "无效的渠道请求头覆盖",
		"ja": "無効なチャンネルヘッダーオーバーライド",
		"fr": "Remplacement d'en-tête de canal invalide",
		"ru": "Недействительное переопределение заголовка канала",
		"vi": "Ghi đè tiêu đề kênh không hợp lệ",
	},
	ErrorCodeChannelModelMappedError: {
		"en": "Channel model mapping error",
		"zh": "渠道模型映射错误",
		"ja": "チャンネルモデルマッピングエラー",
		"fr": "Erreur de mappage de modèle de canal",
		"ru": "Ошибка сопоставления модели канала",
		"vi": "Lỗi ánh xạ mô hình kênh",
	},
	ErrorCodeChannelAwsClientError: {
		"en": "AWS client configuration error",
		"zh": "AWS 客户端配置错误",
		"ja": "AWSクライアント設定エラー",
		"fr": "Erreur de configuration du client AWS",
		"ru": "Ошибка конфигурации клиента AWS",
		"vi": "Lỗi cấu hình client AWS",
	},
	ErrorCodeChannelInvalidKey: {
		"en": "Invalid channel API key",
		"zh": "无效的渠道 API 密钥",
		"ja": "無効なチャンネルAPIキー",
		"fr": "Clé API de canal invalide",
		"ru": "Недействительный ключ API канала",
		"vi": "Khóa API kênh không hợp lệ",
	},
	ErrorCodeChannelResponseTimeExceeded: {
		"en": "Channel response time exceeded",
		"zh": "渠道响应时间超限",
		"ja": "チャンネル応答時間超過",
		"fr": "Temps de réponse du canal dépassé",
		"ru": "Превышено время ответа канала",
		"vi": "Thời gian phản hồi kênh vượt quá giới hạn",
	},
	ErrorCodeChannelNotAvailable: {
		"en": "Channel is not available",
		"zh": "渠道不可用",
		"ja": "チャンネルが利用できません",
		"fr": "Le canal n'est pas disponible",
		"ru": "Канал недоступен",
		"vi": "Kênh không khả dụng",
	},

	// Client Errors (4xxx)
	ErrorCodeReadRequestBodyFailed: {
		"en": "Failed to read request body",
		"zh": "读取请求体失败",
		"ja": "リクエストボディの読み取りに失敗しました",
		"fr": "Échec de la lecture du corps de la requête",
		"ru": "Не удалось прочитать тело запроса",
		"vi": "Không thể đọc nội dung yêu cầu",
	},
	ErrorCodeConvertRequestFailed: {
		"en": "Failed to convert request format",
		"zh": "转换请求格式失败",
		"ja": "リクエストフォーマットの変換に失敗しました",
		"fr": "Échec de la conversion du format de requête",
		"ru": "Не удалось преобразовать формат запроса",
		"vi": "Không thể chuyển đổi định dạng yêu cầu",
	},
	ErrorCodeAccessDenied: {
		"en": "Access denied",
		"zh": "访问被拒绝",
		"ja": "アクセス拒否",
		"fr": "Accès refusé",
		"ru": "Доступ запрещен",
		"vi": "Quyền truy cập bị từ chối",
	},
	ErrorCodeBadRequestBody: {
		"en": "Invalid request body",
		"zh": "无效的请求体",
		"ja": "無効なリクエストボディ",
		"fr": "Corps de requête invalide",
		"ru": "Недействительное тело запроса",
		"vi": "Nội dung yêu cầu không hợp lệ",
	},
	ErrorCodeUnauthorized: {
		"en": "Unauthorized access",
		"zh": "未授权访问",
		"ja": "不正アクセス",
		"fr": "Accès non autorisé",
		"ru": "Неавторизованный доступ",
		"vi": "Truy cập trái phép",
	},
	ErrorCodeForbidden: {
		"en": "Forbidden",
		"zh": "禁止访问",
		"ja": "アクセス禁止",
		"fr": "Interdit",
		"ru": "Запрещено",
		"vi": "Bị cấm",
	},

	// Upstream Errors (5xxx)
	ErrorCodeReadResponseBodyFailed: {
		"en": "Failed to read response body",
		"zh": "读取响应体失败",
		"ja": "レスポンスボディの読み取りに失敗しました",
		"fr": "Échec de la lecture du corps de la réponse",
		"ru": "Не удалось прочитать тело ответа",
		"vi": "Không thể đọc nội dung phản hồi",
	},
	ErrorCodeBadResponseStatusCode: {
		"en": "Bad response status code from upstream",
		"zh": "上游返回错误的状态码",
		"ja": "アップストリームから不正なステータスコードが返されました",
		"fr": "Mauvais code de statut de réponse de l'amont",
		"ru": "Плохой код статуса ответа от восходящего потока",
		"vi": "Mã trạng thái phản hồi không hợp lệ từ phía thượng nguồn",
	},
	ErrorCodeBadResponse: {
		"en": "Bad response from upstream service",
		"zh": "上游服务返回错误响应",
		"ja": "アップストリームサービスから不正な応答がありました",
		"fr": "Mauvaise réponse du service en amont",
		"ru": "Плохой ответ от вышестоящего сервиса",
		"vi": "Phản hồi không hợp lệ từ dịch vụ thượng nguồn",
	},
	ErrorCodeBadResponseBody: {
		"en": "Invalid response body format",
		"zh": "无效的响应体格式",
		"ja": "無効なレスポンスボディフォーマット",
		"fr": "Format de corps de réponse invalide",
		"ru": "Недействительный формат тела ответа",
		"vi": "Định dạng nội dung phản hồi không hợp lệ",
	},
	ErrorCodeEmptyResponse: {
		"en": "Empty response from upstream",
		"zh": "上游返回空响应",
		"ja": "アップストリームからの空の応答",
		"fr": "Réponse vide de l'amont",
		"ru": "Пустой ответ от восходящего потока",
		"vi": "Phản hồi trống từ thượng nguồn",
	},
	ErrorCodeAwsInvokeError: {
		"en": "AWS invocation error",
		"zh": "AWS 调用错误",
		"ja": "AWS呼び出しエラー",
		"fr": "Erreur d'invocation AWS",
		"ru": "Ошибка вызова AWS",
		"vi": "Lỗi gọi AWS",
	},
	ErrorCodeModelNotFound: {
		"en": "Model not found",
		"zh": "未找到模型",
		"ja": "モデルが見つかりません",
		"fr": "Modèle introuvable",
		"ru": "Модель не найдена",
		"vi": "Không tìm thấy mô hình",
	},
	ErrorCodePromptBlocked: {
		"en": "Prompt blocked by content filter",
		"zh": "提示词被内容过滤器阻止",
		"ja": "プロンプトがコンテンツフィルターによってブロックされました",
		"fr": "Invite bloquée par le filtre de contenu",
		"ru": "Подсказка заблокирована контентным фильтром",
		"vi": "Lỗi bị bộ lọc nội dung chặn",
	},
	ErrorCodeRateLimitExceeded: {
		"en": "Rate limit exceeded",
		"zh": "超过速率限制",
		"ja": "レート制限を超過しました",
		"fr": "Limite de taux dépassée",
		"ru": "Превышен лимит скорости",
		"vi": "Vượt quá giới hạn tốc độ",
	},
	ErrorCodeServiceUnavailable: {
		"en": "Service temporarily unavailable",
		"zh": "服务暂时不可用",
		"ja": "サービスは一時的に利用できません",
		"fr": "Service temporairement indisponible",
		"ru": "Сервис временно недоступен",
		"vi": "Dịch vụ tạm thời không khả dụng",
	},

	// Database Errors (6xxx)
	ErrorCodeQueryDataError: {
		"en": "Database query error",
		"zh": "数据库查询错误",
		"ja": "データベースクエリエラー",
		"fr": "Erreur de requête de base de données",
		"ru": "Ошибка запроса к базе данных",
		"vi": "Lỗi truy vấn cơ sở dữ liệu",
	},
	ErrorCodeUpdateDataError: {
		"en": "Database update error",
		"zh": "数据库更新错误",
		"ja": "データベース更新エラー",
		"fr": "Erreur de mise à jour de la base de données",
		"ru": "Ошибка обновления базы данных",
		"vi": "Lỗi cập nhật cơ sở dữ liệu",
	},
	ErrorCodeInsertDataError: {
		"en": "Database insert error",
		"zh": "数据库插入错误",
		"ja": "データベース挿入エラー",
		"fr": "Erreur d'insertion dans la base de données",
		"ru": "Ошибка вставки в базу данных",
		"vi": "Lỗi chèn cơ sở dữ liệu",
	},
	ErrorCodeDeleteDataError: {
		"en": "Database delete error",
		"zh": "数据库删除错误",
		"ja": "データベース削除エラー",
		"fr": "Erreur de suppression de la base de données",
		"ru": "Ошибка удаления из базы данных",
		"vi": "Lỗi xóa cơ sở dữ liệu",
	},
	ErrorCodeDatabaseConnectionFailed: {
		"en": "Database connection failed",
		"zh": "数据库连接失败",
		"ja": "データベース接続に失敗しました",
		"fr": "Échec de la connexion à la base de données",
		"ru": "Не удалось подключиться к базе данных",
		"vi": "Không thể kết nối cơ sở dữ liệu",
	},

	// Quota Errors (7xxx)
	ErrorCodeInsufficientUserQuota: {
		"en": "Insufficient user quota",
		"zh": "用户配额不足",
		"ja": "ユーザークォータが不足しています",
		"fr": "Quota utilisateur insuffisant",
		"ru": "Недостаточная квота пользователя",
		"vi": "Hạn ngạch người dùng không đủ",
	},
	ErrorCodePreConsumeTokenQuotaFailed: {
		"en": "Failed to pre-consume token quota",
		"zh": "预消耗 token 配额失败",
		"ja": "トークンクォータの事前消費に失敗しました",
		"fr": "Échec de la pré-consommation du quota de jetons",
		"ru": "Не удалось предварительно израсходовать квоту токенов",
		"vi": "Không thể tiêu thụ hạn ngạch token trước",
	},
	ErrorCodeQuotaExceeded: {
		"en": "User quota exceeded",
		"zh": "超出用户配额",
		"ja": "ユーザークォータを超過しました",
		"fr": "Quota utilisateur dépassé",
		"ru": "Превышена квота пользователя",
		"vi": "Vượt quá hạn ngạch người dùng",
	},
}

// Localize returns localized error message based on the language code
// Falls back to English if the requested language is not available
func (e *NewAPIError) Localize(lang string) string {
	if e == nil {
		return ""
	}

	msgs, ok := errorMessages[e.errorCode]
	if !ok {
		// Fallback to the error message from Err
		return e.Error()
	}

	// Try to get message in requested language
	if msg, ok := msgs[lang]; ok {
		return msg
	}

	// Fallback to English
	if msg, ok := msgs["en"]; ok {
		return msg
	}

	// Final fallback to the error message from Err
	return e.Error()
}

// GetLanguageFromContext extracts language code from Accept-Language header
// Supports formats like "zh-CN", "en-US", "ja", etc.
// Returns the language code (e.g., "zh", "en", "ja")
func GetLanguageFromContext(acceptLanguage string) string {
	if acceptLanguage == "" {
		return "en" // default to English
	}

	// Parse language header (e.g., "zh-CN" -> "zh")
	parts := strings.Split(acceptLanguage, "-")
	if len(parts) > 0 {
		lang := strings.ToLower(parts[0])
		// Validate that we support this language
		if isLanguageSupported(lang) {
			return lang
		}
	}

	return "en" // fallback to English
}

// isLanguageSupported checks if the language code is supported
func isLanguageSupported(lang string) bool {
	supportedLanguages := map[string]bool{
		"en": true,
		"zh": true,
		"ja": true,
		"fr": true,
		"ru": true,
		"vi": true,
	}
	return supportedLanguages[lang]
}

// GetSupportedLanguages returns a list of supported language codes
func GetSupportedLanguages() []string {
	return []string{"en", "zh", "ja", "fr", "ru", "vi"}
}

// GetLanguageName returns the full name of a language code
func GetLanguageName(lang string) string {
	languageNames := map[string]string{
		"en": "English",
		"zh": "中文",
		"ja": "日本語",
		"fr": "Français",
		"ru": "Русский",
		"vi": "Tiếng Việt",
	}
	if name, ok := languageNames[lang]; ok {
		return name
	}
	return lang
}
