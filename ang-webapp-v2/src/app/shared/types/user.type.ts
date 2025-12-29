import { Audit } from "./common.type";

export interface AppCustomer {
  nimb_id: string;
  fname: string;
  lname: string;
  email: string;
  password: string;
  organization: Organization;
  role?: Role;
  audit?: Audit;
}

export interface Organization {
  nimb_id: string;
  name: string;
  address: string;
  audit: Audit;
}

export interface Role {
  name: string;
  permissions: string[];
}
