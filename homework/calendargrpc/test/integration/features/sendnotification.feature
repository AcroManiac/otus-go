
Feature: Send notification to event's owner
  As owner of event created beforehand
  In order to understand that the notification was sent
  I want to read information from PostgreSQL service

  Scenario: Read information from database table
    Given Connection to PostgreSQL service with DSN "postgres://dbuser:En9NR2b869@postgres:5432/calendar?sslmode=disable"
    When I wait for 1 minute
    Then selection from table "notices" should return 1 notice