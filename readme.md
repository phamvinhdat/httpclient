# Simple http client for go, easy to use

## Example

```golang
result := Credential{}
client := httpclient.NewClient(
  	httpclient.WithSender(gosender.New(gosender.WithTimeout(time.Second * 5))),
)
statusCode, err := client.Post(context.Background(), "http://localhost:8080/login",
   	httpclient.WithBodyProvider(body.NewJson(Credential{
   		Username: "admin",
   		Password: "admin",
   	})),
   	httpclient.WithHookFn(hook.UnmarshalResponse(&result)),
   	httpclient.WithHookFn(hook.Log()),
)
```
