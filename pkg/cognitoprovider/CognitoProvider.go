package cognitoprovider 

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
    "github.com/aws/aws-sdk-go/service/cognitoidentity"
    "time"
    "errors"
)

type cognitoProvider struct {
    Region *string
    Username *string
    ClientId *string
    UserPoolId *string
    IdentityPoolId *string
    IdToken *string
    Credentials *cognitoidentity.Credentials
}

func New(region, username, clientId, userPoolId, identityPoolId *string) credentials.Provider {
    return &cognitoProvider{
        Region: region,
        Username: username,
        ClientId: clientId,
        UserPoolId: userPoolId,
        IdentityPoolId: identityPoolId,
    }
}

func (cp *cognitoProvider) Retrieve() (credentials.Value, error) {
    if cp.Credentials == nil {
        token, err := cp.InitiateAuth()

        if err != nil {
            return credentials.Value{}, err
        }

        creds, err := cp.GetCredentialsForIdentity(token)

        if err != nil {
            return credentials.Value{}, err
        }

        cp.Credentials = creds
    } 
    
    return credentials.Value {
        AccessKeyID: *cp.Credentials.AccessKeyId,
        SecretAccessKey: *cp.Credentials.SecretKey,
        SessionToken: *cp.Credentials.SessionToken,
        ProviderName: "CognitoProvider",
    }, nil
}

func (cp *cognitoProvider) IsExpired() bool {
    if cp.Credentials == nil {
        return true
    } else {
        return cp.Credentials.Expiration.After(time.Now())
    }
}

func (cp *cognitoProvider) InitiateAuth() (*string, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(*cp.Region),
    })

    if err != nil {
        return nil, err
    }
   
    fmt.Print("Please Enter A Password: ")

    var secret string
    _, err = fmt.Scan(&secret)
    if err != nil {
        return nil, err
    }

    svc := cognitoidentityprovider.New(sess)
    resp, err := svc.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
        AuthFlow: aws.String("USER_PASSWORD_AUTH"),
        AuthParameters: map[string]*string {
            "USERNAME": cp.Username,
            "PASSWORD": aws.String(string(secret)),
        },
        ClientId: cp.ClientId,
    })

    if err != nil {
        return nil, err
    } else if resp.ChallengeName != nil && *resp.ChallengeName == "NEW_PASSWORD_REQUIRED" {
        return cp.NewPasswordChallengeResponse(resp.Session)
    } else if resp.ChallengeName != nil {
        return nil, errors.New(fmt.Sprintf("An unsupported authentication challenge was encountered: %s", *resp.ChallengeName))
    } else {
        return resp.AuthenticationResult.IdToken, nil
    }
}

func (cp *cognitoProvider) NewPasswordChallengeResponse(sessionId *string) (*string, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(*cp.Region),
    })

    if err != nil {
        return nil, err
    }
    
    fmt.Println("Your password needs to be changed.")
    fmt.Print("Please re-enter your old password: ")
    var oldPassword string
    _, err = fmt.Scan(&oldPassword)
    if err != nil {
        return nil, err
    }
    fmt.Print("Please enter a new password: ")
    var newPassword string
    _, err = fmt.Scan(&newPassword)
    if err != nil {
        return nil, err
    }

    svc := cognitoidentityprovider.New(sess)
    resp, err := svc.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
        ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
        ChallengeResponses: map[string]*string {
            "USERNAME": cp.Username,
            "PASSWORD": aws.String(string(oldPassword)),
            "NEW_PASSWORD": aws.String(string(newPassword)),
        },
        ClientId: cp.ClientId,
        Session: sessionId,
    })

    if err != nil {
        return nil, err
    } else {
        return resp.AuthenticationResult.IdToken, nil
    }
}

func (cp *cognitoProvider) GetCredentialsForIdentity(idToken *string) (*cognitoidentity.Credentials, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: cp.Region,
    })

    if err != nil {
        return nil, err
    }

    svc := cognitoidentity.New(sess)

    idresp, err := svc.GetId(&cognitoidentity.GetIdInput{
        IdentityPoolId: cp.IdentityPoolId,
        Logins: map[string]*string {
            fmt.Sprintf("cognito-idp.%s.amazonaws.com/%s", *cp.Region, *cp.UserPoolId): idToken,
        },
    })

    if err != nil {
        return nil, err
    }

    resp, err := svc.GetCredentialsForIdentity(&cognitoidentity.GetCredentialsForIdentityInput{
        IdentityId: idresp.IdentityId,
        Logins: map[string]*string {
            fmt.Sprintf("cognito-idp.%s.amazonaws.com/%s", *cp.Region, *cp.UserPoolId): idToken,
        },
    })

    if err != nil {
        return nil, err
    } else {
        return resp.Credentials, nil
    }
}
