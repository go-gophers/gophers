-- local base = "https://api.github.com/?access_token=" .. os.getenv("GOPHERS_GITHUB_TOKEN")
-- local client = gophers.new_client(base)

local gophers = require("gophers")

local response, err = gophers.get("https://httpbin.org/headers")

print("Status:", response.status)
print("Code:", response.code)

print("Body:")
print(response.body)
print("Body length:", string.len(response.body))

print("Error:", err)
