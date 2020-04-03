
Feature: Get calendar events for time period
  As API client of Calendar service
  In order to understand that the storage contains Calendar events
  I want to get events for specified time period from Calendar API service

  Scenario: Calendar event is obtained
    Given Connection to Calendar API on "127.0.0.1:8888"
    When I send GetEvents request with period "day" and start time "2020-04-02T12:00:00+03:00"
    Then search should return 1 event