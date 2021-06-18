# Secret Manager 관리 API 

AWS Secretsmanager 를 이용하는 코드를 작성한다. 

## 의존성 라이브러리 다운로드 

```go
go get "github.com/aws/aws-sdk-go/aws"
go get "github.com/aws/aws-sdk-go/aws/awserr"
go get "github.com/aws/aws-sdk-go/aws/session"
go get "github.com/aws/aws-sdk-go/service/secretsmanager"
```

## 프로그램 기본 구조 

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	PROFILE = "sdk-user"
	REGION  = "ap-northeast-2"
)

var (
	SecretName  = "myproject/devops/configs"
	Description = "This secret is my Project Secret"
	Secrets     = map[string]string{
		"myDBName":     "dbname",
		"myDBPassword": "dbpassword",
		"dbPort":       "dbPorts",
	}
	UpdateSecrets = map[string]string{
		"myDBName":     "modDbname",
		"myDBPassword": "modDbpassword",
		"dbPort":       "modDbPorts",
	}
	PutSecrets = map[string]string{
		"myDBName":     "putDbname",
		"myDBPassword": "putDbpassword",
		"dbPort":       "putDbPorts",
	}
)

// getSecretManager()
// secret manager 에 접근하기 위한 세션을 생성한다.
func getSecretManager() *secretsmanager.SecretsManager {...}

// Secret 을 신규로 생성한다.
func createSecrets(smgr *secretsmanager.SecretsManager, inputValue *secretsmanager.CreateSecretInput) *secretsmanager.CreateSecretOutput {...}

// Secret 정보를 조회한다.
func getSecretInfos(smgr *secretsmanager.SecretsManager, secretId string) *secretsmanager.DescribeSecretOutput {
	input := &secretsmanager.DescribeSecretInput{...}

// 시크릿 키와 값을 조회한다.
func getSecretValue(smgr *secretsmanager.SecretsManager, secretId string) *secretsmanager.GetSecretValueOutput {...}

// 시크릿 정보를 업데이트한다.
// 기본 시크릿 정보를 변경하고,
func updateSecret(smgr *secretsmanager.SecretsManager, input *secretsmanager.UpdateSecretInput) *secretsmanager.UpdateSecretOutput {...}

// 시크릿 정보를 업데이트한다.
func putSecretValue(smgr *secretsmanager.SecretsManager, input *secretsmanager.PutSecretValueInput) *secretsmanager.PutSecretValueOutput {...}

func deleteSecret(smgr *secretsmanager.SecretsManager, secretId string, optional ...int64) *secretsmanager.DeleteSecretOutput {...}

func main() {

	smgr := getSecretManager()

	secretString, _ := json.Marshal(Secrets)

	input := &secretsmanager.CreateSecretInput{
		Description:  aws.String(Description),
		Name:         aws.String(SecretName),
		SecretString: aws.String(string(secretString)),
	}

	result := createSecrets(smgr, input)
	fmt.Println("----- Create Secret Success: ----- ")
	fmt.Println(result)

	secretInfo := getSecretInfos(smgr, SecretName)
	fmt.Println("----- Get Secret Infos: ----- ")
	fmt.Println(secretInfo)

	secretPolicy := getResourcePolicy(smgr, SecretName)
	fmt.Println("----- Get Resource Policy: ----- ")
	fmt.Println(secretPolicy)

	myDBPassword := getSecretValue(smgr, SecretName)
	fmt.Println("----- Get My DB Password: ----- ")
	fmt.Println(myDBPassword)

	modSecretString, _ := json.Marshal(UpdateSecrets)
	updateValues := &secretsmanager.UpdateSecretInput{
		Description:  aws.String("Modified Descriptions."),
		SecretId:     aws.String(SecretName),
		SecretString: aws.String(string(modSecretString)),
	}

	updateResult := updateSecret(smgr, updateValues)
	fmt.Println("----- Update Secrets: ----- ")
	fmt.Println(updateResult)

	putSecretString, _ := json.Marshal(PutSecrets)
	putValues := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(SecretName),
		SecretString: aws.String(string(putSecretString)),
	}

	putResult := putSecretValue(smgr, putValues)
	fmt.Println("----- Put Secrets: ----- ")
	fmt.Println(putResult)

	deleteResult := deleteSecret(smgr, SecretName, 30)
	fmt.Println("----- Delete Secrets: ----- ")
	fmt.Println(deleteResult)

}

```

위 항목과 같이 SecretsManager 를 운용하기 위한 소스 기본 구조를 정의 했다. 

- getSecretManager: 시크릿 매니저 핸들러를 생성하고, 획득한다. 
- createSecrets: 신규 시크릿을 생성한다. 
- getSecretInfos: 시크릿 정보를 조회한다. 
- getResourcePolicy: 시크릿 리소스의 운영 정책을 조회한다.
- getSecretValue: 시크릿 키에 대한 값을 조회한다. 
- updateSecret: 시크릿 정보를 업데이트한다. 
- putSecretValue: 시크릿 정보를 업데이트하고, 새로운 키/값을 추가한다.
- deleteSecret: 시크릿 정보를 삭제한다. (시크릿 정보는 보통 특정 기간을 두고 삭제 스케줄링을 걸게 되어 있다. )
  - 매우 민감한 데이터이므로 함부로 삭제되지 않고, 스케줄을 두어 삭제하도록 한다. 

## Secret 매니저 생성하기 

getSecretManager 메소드 구현을 살펴 보자. 

```go
func getSecretManager() *secretsmanager.SecretsManager {
	// Secret Manager 에 접근할 세션을 생성한다.
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Profile: "sdk-user",
			Config: aws.Config{
				Region: aws.String("ap-northeast-2"),
			},
		},
	)

	if err != nil {
		log.Fatalf("Cannot create session %v", err)
	}

	smgr := secretsmanager.New(sess)

	return smgr
}
```

aws config 정보에서 특정 프로파일을 조회하고, secretsmanager 의 리젼을 설정한다.

이렇게 설정된 정보를 기준으로 secretsmanager 를 생성하여 반환한다. 

## 시크릿 생성하기. 

시크릿은 특정 서비스에 사용할 시크릿 정보를 생성하는 것이다. 

createSecrets 메소드로 구현하며, 시크릿 키, 설명, 시크릿 내용등을 기술하여 생성할 수 있다. 

```go
func createSecrets(smgr *secretsmanager.SecretsManager, inputValue *secretsmanager.CreateSecretInput) *secretsmanager.CreateSecretOutput {
	result, err := smgr.CreateSecret(inputValue)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodePreconditionNotMetException:
				fmt.Println(secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}
```

시크릿을 생성하는 코드이다. 

파라미터로 시크릿의 정보를 담고 있는 생성 파라미터 객체에 담아서 생성하게 된다. 

```go
	secretString, _ := json.Marshal(Secrets)

	input := &secretsmanager.CreateSecretInput{
		Description:  aws.String(Description),
		Name:         aws.String(SecretName),
		SecretString: aws.String(string(secretString)),
	}
```

CreateSecretInput 이 해당 요청 파라미터 객체이다. 

### 결과 보기

```go
----- Create Secret Success: ----- 
{
  ARN: "arn:aws:secretsmanager:ap-northeast-2:103382364946:secret:myproject/devops/configs/v1-37mG8J",
  Name: "myproject/devops/configs/v1",
  VersionId: "9B61A56B-96F8-45A1-8838-3F3BC46FE330"
}
```

## 시크릿 자체 정보를 조회한다. 

getSecretInfos 메소드 구현은 다음과 같다. 

```go
func getSecretInfos(smgr *secretsmanager.SecretsManager, secretId string) *secretsmanager.DescribeSecretOutput {
	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretId),
	}

	result, err := smgr.DescribeSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}
```

시크릿을 설명하는 정보를 조회할 수 있다. 

### 결과

```go
----- Get Secret Infos: ----- 
{
  ARN: "arn:aws:secretsmanager:ap-northeast-2:103382364946:secret:myproject/devops/configs/v1-37mG8J",
  CreatedDate: 2021-03-17 07:08:30.469 +0000 UTC,
  Description: "This secret is my Project Secret",
  LastChangedDate: 2021-03-17 07:08:30.518 +0000 UTC,
  Name: "myproject/devops/configs/v1",
  VersionIdsToStages: {
    9B61A56B-96F8-45A1-8838-3F3BC46FE330: ["AWSCURRENT"]
  }
}
```

## 시크릿 운영 정책 정보를 조회한다. 

getResourcePolicy 메소드 구현은 다음과 같다. 

```go
func getResourcePolicy(smgr *secretsmanager.SecretsManager, secretId string) *secretsmanager.GetResourcePolicyOutput {
	input := &secretsmanager.GetResourcePolicyInput{
		SecretId: aws.String(secretId),
	}

	result, err := smgr.GetResourcePolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}
```

### 결과 보기 

```go
----- Get Resource Policy: ----- 
{
  ARN: "arn:aws:secretsmanager:ap-northeast-2:103382364946:secret:myproject/devops/configs/v1-37mG8J",
  Name: "myproject/devops/configs/v1"
}
```

## 저장된 시크릿 키와 값을 조회한다. 

getSecretValue 코드는 다음과 같다. 

```go
func getSecretValue(smgr *secretsmanager.SecretsManager, secretId string) *secretsmanager.GetSecretValueOutput {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}

	result, err := smgr.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result

}
```

### 결과

```go
----- Get My DB Password: ----- 
{
  ARN: "arn:aws:secretsmanager:ap-northeast-2:103382364946:secret:myproject/devops/configs/v1-37mG8J",
  CreatedDate: 2021-03-17 07:08:30.513 +0000 UTC,
  Name: "myproject/devops/configs/v1",
  SecretString: "{\"dbPort\":\"dbPorts\",\"myDBName\":\"dbname\",\"myDBPassword\":\"dbpassword\"}",
  VersionId: "9B61A56B-96F8-45A1-8838-3F3BC46FE330",
  VersionStages: ["AWSCURRENT"]
}
```

## 기존 시크릿 정보 업데이트하기. 

updateSecret 코드는 다음과 같다. 

```go
func updateSecret(smgr *secretsmanager.SecretsManager, input *secretsmanager.UpdateSecretInput) *secretsmanager.UpdateSecretOutput {

	result, err := smgr.UpdateSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				fmt.Println(secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodePreconditionNotMetException:
				fmt.Println(secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}
```

### 결과

```go
----- Update Secrets: ----- 
{
  ARN: "arn:aws:secretsmanager:ap-northeast-2:103382364946:secret:myproject/devops/configs/v1-37mG8J",
  Name: "myproject/devops/configs/v1",
  VersionId: "1CC966F8-D9EC-42F6-B068-43E1B2E44730"
}
```

## 시크릿 정보를 업데이트 하고, 새로운 키/값을 추가한다. 

putSecretValue 코드는 다음과 같다. 

```go
func putSecretValue(smgr *secretsmanager.SecretsManager, input *secretsmanager.PutSecretValueInput) *secretsmanager.PutSecretValueOutput {

	result, err := smgr.PutSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}
```

### 결과 

```go
----- Put Secrets: ----- 
{
  ARN: "arn:aws:secretsmanager:ap-northeast-2:103382364946:secret:myproject/devops/configs/v1-37mG8J",
  Name: "myproject/devops/configs/v1",
  VersionId: "F7F37D5D-65BE-4094-8CDF-CC0D85C47D48",
  VersionStages: ["AWSCURRENT"]
}
```

## 시크릿 삭제하기. 

deleteSecret 코드는 다음과 같다. 

```go
func deleteSecret(smgr *secretsmanager.SecretsManager, secretId string, optional ...int64) *secretsmanager.DeleteSecretOutput {
	recoveryDay := int64(7)
	if len(optional) > 0 {
		recoveryDay = optional[0]
	}
	input := &secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(recoveryDay),
		SecretId:             aws.String(secretId),
	}

	result, err := smgr.DeleteSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}
```

### 결과

```go
----- Delete Secrets: ----- 
{
  ARN: "arn:aws:secretsmanager:ap-northeast-2:103382364946:secret:myproject/devops/configs/v1-37mG8J",
  DeletionDate: 2021-04-16 07:08:30.763 +0000 UTC,
  Name: "myproject/devops/configs/v1"
}
```

## 결론

지금까지 SecretManager 을 관리하는 SDK 를 사용하는 방법을 살펴 보았다.

SecretManager 를 위해서 AWS Console 에 접근하여, 키/값을 추가할 수 있다. 그러나 다양한 프로젝트가 운영되는 환경이라면 위 코드를 활용하여 각 프로젝트별 secret 을 공통관리 할 수 있는 방향으로 개발을 진행할 수 있을 것이다. 

참고:

- aws sdk library: https://aws.github.io/aws-sdk-go-v2/docs/code-examples/
  