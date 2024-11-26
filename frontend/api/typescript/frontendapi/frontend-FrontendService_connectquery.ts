// @generated by protoc-gen-connect-query v2.0.0 with parameter "target=ts"
// @generated from file frontendapi/frontend.proto (package frontendapi, syntax proto3)
/* eslint-disable */

import { FrontendService } from "./frontend_pb";

/**
 * Gets information for the current user.
 *
 * @generated from rpc frontendapi.FrontendService.GetUser
 */
export const getUser = FrontendService.method.getUser;

/**
 * Saves information for a user. This method works both for a new or existing user.
 * The user is identified by the firebase ID token included in the authorization header.
 *
 * @generated from rpc frontendapi.FrontendService.SaveUser
 */
export const saveUser = FrontendService.method.saveUser;
