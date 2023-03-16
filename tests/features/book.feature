Feature: Book management
    In order to use book API
    As a Librarian
    I need to be able to manage books

    Scenario: then user try to insert one book, one created book should be displayed by the system
        When I send "POST" request to "/books" with payload:
        """
        {
            "id": 1,
            "title": "Dune",
            "author": "Frank Herbert"
        }   
        """
        Then the response code should be 201
        And the response payload should match json:
        """
        [
            {
                "id": 1,
                "title": "Dune",
                "author": "Frank Herbert"
            }
        ]   
        """