Feature: User Authentication
  Scenarios for user registration, login, password change, and login with new password.

  Scenario: Register and login user
    Given a user with the following details:
      | ID        | Fullname | Login  | Password | Birthdate   | Email                      | Phone           |
      | 9cb8f50b-66c7-4dbd-a1c6-a589e1471360 | sacha   | sacha  | 12345    | 1957-10-10 | sachatarba@rambler.ru | +7-985-985-99-99 |
    When the user registers
    Then the status code should be "200 OK"
    When the user logs in
    Then the status code should be "200 OK"
    When the confirmation code is fetched from the email "sachatarba@rambler.ru"
    And the user confirms the account
    Then the status code should be "200 OK"

  Scenario: Change password and login with new password
    Given a user with the following details:
      | ID        | Fullname | Login   | Password | Birthdate   | Email                      | Phone           |
      | 9cb8f50b-66c7-4dbd-a1c6-a589e1471361 | sacha1  | sacha1 | 12345    | 1957-10-10 | sachatarba@rambler.ru | +7-985-985-99-99 |
    When the user registers
    Then the status code should be "200 OK"
    When the user logs in
    Then the status code should be "200 OK"
    When the confirmation code is fetched from the email "sachatarba@rambler.ru"
    And the user confirms the account
    Then the status code should be "200 OK"
    When the user changes the password to "123456"
    Then the status code should be "200 OK"
    When the user logs in
    Then the status code should be "200 OK"
    When the confirmation code is fetched from the email "sachatarba@rambler.ru"
    And the user confirms the account
    Then the status code should be "200 OK"
