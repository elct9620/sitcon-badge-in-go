package main

import(
  "github.com/codegangsta/martini"
  "net/http"
  "time"
  "encoding/json"
  "encoding/hex"
  "crypto/sha256"
)

// Data Struct
type BadgeAssertion struct {
  Uid string `json:"uid"`
  Recipient IdentityObject `json:"recipient"`
  Badge string `json:"badge"`
  Verify VerificationObject `json:"verify"`
  IssuedOn int64 `json:"issuedOn"`
}

type IdentityObject struct {
  Identity string `json:"identity"`
  Type string `json:"type"`
  Hashed bool `json:"hashed"`
  Salt string `json:"salt"`
}

type VerificationObject struct {
  Type string `json:"type"`
  URL string `json:"url"`
}

type Badge struct {
  Name string `json:"name"`
  Description string `json:"description"`
  Image string `json:"image"`
  Criteria string `json:"criteria"`
  Issuer string `json:"issuer"`
}

type Issuer struct {
  Name string `json:"name"`
  URL string `json:"url"`
}

// Middware
func jsonMiddleware(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "application/json")
}

// Handlers
func handleBadge(req *http.Request) string {
  badge := &Badge{
    Name: "Test",
    Description: "A test badge",
    Image: "http://" + req.Host + "/images/badge.png",
    Criteria: "http://" + req.Host + "/badge.html",
    Issuer: "http://" + req.Host + "/organization",
  }
  badgeJSON, err := json.Marshal(badge)
  if err != nil {
    return "{\"error\": \"JSON Parse error\"}"
  }
  return string(badgeJSON)
}

func handleBadgeDetial(res http.ResponseWriter) string {
  res.Header().Set("Content-Type", "text/html")
  return "Just for test..."
}

func handleAssertion(req *http.Request) string {
  hash := sha256.New()
  hash.Write([]byte("example@sitcon.org"))
  md := hash.Sum([]byte("salt"))
  identity := &IdentityObject{
    Identity: "sha256$" + hex.EncodeToString(md),
    Type: "email",
    Hashed: true,
    Salt: "salt",
  }
  verify := &VerificationObject{Type: "hosted", URL: "http://" + req.Host + "/assertion"}
  issued := time.Now()
  assertion := &BadgeAssertion{
    Uid: "badge-1",
    Recipient: *identity,
    Badge: "http://" + req.Host + "/badge",
    Verify: *verify,
    IssuedOn: issued.Unix(),
  }
  assertionJSON, _ := json.Marshal(assertion)
  return string(assertionJSON)
}

func handleIssuer() string {
  issuer := &Issuer{Name: "Frost Studio", URL: "http://frost.tw"}
  issuerJSON, _ := json.Marshal(issuer)
  return string(issuerJSON)
}

func main() {
  m := martini.Classic()
  m.Use(martini.Static("assets"))
  m.Use(jsonMiddleware)
  m.Get("/", func () string {
    return "Hello World"
  })
  m.Get("/badge", handleBadge)
  m.Get("/badge.html", handleBadgeDetial)
  m.Get("/assertion", handleAssertion)
  m.Get("/organization", handleIssuer)
  m.Run()
}
