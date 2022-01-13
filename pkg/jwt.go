package pkg

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var jwtSecret []byte
var jwtExpireHour int

func SetUpJwt() {
	var jwtConfig = viper.GetStringMap("jwt")
	jwtSecret = []byte(jwtConfig["secret"].(string))
	jwtExpireHour = jwtConfig["expire_hour"].(int)
}

type Claims struct {
	UserId uint   `json:"userid"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err

}

func GenerateToken(userid uint, email string) (string, error) {
	log.Println("------usrid: ", userid, email)
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(jwtExpireHour) * time.Hour)

	claims := Claims{
		userid,
		email,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "workerunion",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("secret: ", jwtSecret)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}
