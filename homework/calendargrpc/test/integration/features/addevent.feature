
Feature: Add calendar event
  As API client of Calendar service
  In order to understand that the user created Calendar event
  I want to receive response with added event from Calendar API service

  Scenario: Calendar event is added
    Given Connection to Calendar API on "127.0.0.1:8888"
    And There is the event:
    """
		{
			"title": "Event 1",
			"description": "Data for testing microservices",
			"owner": "Artem",
			"startTime": "2020-04-02T12:03:00+03:00"
		}
	"""
    When I send AddEvent request
    Then response should have event