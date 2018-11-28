CREATE TABLE USER_ROLE_MAP(
  ROLE_ID INT(6),
  UID VARCHAR(36),
  PRIMARY KEY(UID, ROLE_ID),
  FOREIGN KEY(ROLE_ID) REFERENCES ROLES_REF(ID),
  FOREIGN KEY(UID) REFERENCES USERS(UID)
);
