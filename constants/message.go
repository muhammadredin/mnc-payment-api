package constants

// This file used to store response message as a constant
const CustomerCreateSuccess = "Successfully created a new customer"
const CustomerCreateError = "An error occurred while creating the customer"
const CustomerNotFound = "Customer not found"
const CustomerFindSuccess = "Successfully get a customer"
const UsernameDuplicateError = "Username already exists"

const RefreshTokenNotFoundError = "Refresh token not found, could be expired"
const RefreshTokenExpiredError = "Refresh token is expired"

const LoginSuccess = "Successfully logged in"
const LoginUnauthorizedError = "Invalid credentials"
const AuthorizationHeaderMissingError = "Authorization header is missing"
const AuthorizationHeaderInvalidError = "Invalid Authorization header token"
const LogoutSuccess = "Successfully logged out"
const AccessTokenNotFoundError = "Access token not found"
const AuthenticatedUserNotFoundError = "Authenticated user not found"

const JsonWriteSuccess = "Successfully wrote JSON file"
const JsonFileNotFound = "Json file path not found"
const JsonMarshalError = "An error occurred while marshalling json"
const JsonMappingError = "An error occurred while mapping json, make sure the data type is compatible with the json"
const JsonCreateError = "An error occurred while creating JSON file"
const JsonWriteError = "An error occurred while writing JSON file"

const JwtTokenInvalidError = "Invalid JWT token"

const WalletNotFoundError = "Wallet not found"
const WalletDuplicateError = "Wallet already exists"
const WalletForbiddenAccess = "User does not have permission to access this wallet"

const InvalidRequestBodyError = "Invalid request body"

const TransactionInsufficientError = "Insufficient amount of funds"
const TransactionSuccess = "Successfully created a transaction"
