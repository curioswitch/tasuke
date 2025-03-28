
�
frontendapi/frontend.protofrontendapi"j
User8
programming_language_ids (RprogrammingLanguageIds(
max_open_reviews (RmaxOpenReviews"
GetUserRequest"8
GetUserResponse%
user (2.frontendapi.UserRuser"8
SaveUserRequest%
user (2.frontendapi.UserRuser"
SaveUserResponse2�
FrontendServiceD
GetUser.frontendapi.GetUserRequest.frontendapi.GetUserResponseG
SaveUser.frontendapi.SaveUserRequest.frontendapi.SaveUserResponseB;Z9github.com/curioswitch/tasuke/frontend/api/go;frontendapiJ�

  1

  

 

 P
	
 P
&
   The settings for a user.



 
�
  /� IDs of programming languages that reviews can be created for.
 IDs correspond to `language_id` from github-linguist.
 https://github.com/github-linguist/linguist/blob/master/lib/linguist/languages.yml
 Required.


  


  

  *

  -.
d
 W The maximum number of reviews created by the app that can be open at once.
 Required.


 

 	

 
3
 ( A request for FrontendService.GetUser.




5
 ) A response for FrontendService.GetUser.




/
 " The user information.
 Required.


 

 

 
5
 ") A request for FrontendService.SaveUser.




-
 !  The user to create.
 Required.


 !

 !

 !
Z
% '* A response for FrontendService.SaveUser.
"" Empty to allow future extension.



%
+
 * 1 The service for the frontend.



 *
5
  ,8( Gets information for the current user.


  ,

  ,

  ,'6
�
 0;� Saves information for a user. This method works both for a new or existing user.
 The user is identified by the firebase ID token included in the authorization header.


 0

 0

 0)9bproto3