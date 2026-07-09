package main 
import ("fmt" ; "net/http" ; "encoding/json" ; "time" ; "context" ; "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" ; "go.mongodb.org/mongo-driver/bson/primitive" ; "github.com/golang-jwt/jwt" ; "strings" )

var client *mongo.Client
var apicollection *mongo.Collection
var checkcollection *mongo.Collection
var usercollection *mongo.Collection



type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username         string             `bson:"username" json:"username"`
	email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"password"`
}


type API struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"userId" json:"userId"`

	Name             string             `bson:"name" json:"name"`
	URL              string             `bson:"url" json:"url"`

	LastStatus       string             `bson:"lastStatus" json:"lastStatus"`
	LastStatusCode   int                `bson:"lastStatusCode" json:"lastStatusCode"`
	LastResponseTime int64              `bson:"lastResponseTime" json:"lastResponseTime"`
	LastCheckedAt    time.Time          `bson:"lastCheckedAt" json:"lastCheckedAt"`

	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
}


type Check struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	APIID            primitive.ObjectID `bson:"apiId" json:"apiId"`

	Status           string             `bson:"status" json:"status"`
	StatusCode       int                `bson:"statusCode" json:"statusCode"`
	ResponseTime     int64              `bson:"responseTime" json:"responseTime"`

	CheckedAt        time.Time          `bson:"checkedAt" json:"checkedAt"`
}

type loginuser struct {
	email   string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ;
	}
	var user User;
	err := json.NewDecoder(r.Body).Decode(&user);
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest);
		return ;
	}


	count , err := usercollection.CountDocuments(context.Background(), bson.M{"username": user.Username , "email": user.email , "password": user.Password});
	if err != nil {
		http.Error(w, "Failed to check existing user", http.StatusInternalServerError)
		return ;
	}
	
	if count > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return ;
	}
	result, err := usercollection.InsertOne(context.Background(), user);
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError);
		return ;
	}

	fmt.Println("User created with ID: ", result.InsertedID);

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ;
	}
	var loginUser loginuser;
	err := json.NewDecoder(r.Body).Decode(&loginUser);
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return ;
	}

	var user User;
	err = usercollection.FindOne(context.Background(), bson.M{"email": loginUser.email}).Decode(&user);

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return ;
	}

	if user.Password != loginUser.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return ;
	}


	claims := jwt.MapClaims{
		"userId": user.ID ,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return ;
	}


	loginResp := loginResponse{
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(loginResp);	


}



func main() {
	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second);
	defer cancel();
	var err error;
	client , err = mongo.Connect(ctx , options.Client().ApplyURI("mongodb://localhost:27017"));
	if err != nil {
		fmt.Println("Error connecting to MongoDB: " , err);
		return ;
	}

	apicollection = client.Database("amd").Collection("apicollection");
	checkcollection = client.Database("amd").Collection("checkcollection");
	usercollection = client.Database("amd").Collection("usercollection");
	

	http.HandleFunc("/signup", enableCORS(createUser));
	http.HandleFunc("/login", enableCORS(login));
	
	fmt.Println("Server running on http://localhost:8080");
	 http.ListenAndServe(":8080", nil);

}