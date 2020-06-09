package app_routes

const Home = ""
const HomeSlash = "/"

const Faq = "/faq"
const About = "/about"

const Login = "/login"
const CreateUser = "/createUser"
const ForgotUsernameOrPassword = "/forgotUsernameOrPassword"
const Logout = "/logout"

const ActivateAccountBase = "/activateAccount/"
const ActivateAccount = ActivateAccountBase + "{rediskey}"

const ResetPasswordBase = "/resetPassword/"
const ResetPassword = ResetPasswordBase + "{rediskey}"

const AddNewWebsite = "/addNewWebsite"
const GenerateNewPassword = "/generateNewPassword"
const ViewPassword = "/viewPassword"

const AddNewSecret = "/addNewSecret"
const ViewSecret = "/viewSecret"
const DeleteSecret = "/deleteSecret"
