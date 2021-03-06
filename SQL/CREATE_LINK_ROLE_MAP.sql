CREATE TABLE LINK_ROLE_MAP (
  ROLE_ID INT(6),
  LINK_ID INT(6),
  PRIMARY KEY(ROLE_ID, LINK_ID),
  FOREIGN KEY(ROLE_ID) REFERENCES ROLES_REF(ID),
  FOREIGN KEY(LINK_ID) REFERENCES LINKS_REF(ID)
);
