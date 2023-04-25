import { uid } from 'uid';

describe('Create message', () => {
    
    it('Check if create message box exists', () => {

        // ---------------- SETUP ---------------- //

        const username = uid()

        // register
        cy.register(username, "user@example.com", "password", "password")

        // login
        cy.login(username, "password")
        
        // visit user page
        cy.visit('localhost:8080/'+username)

        // ---------------- TEST ---------------- //

        // check if page contains 'What's on your mind ?', share button and input text field
        cy.contains('What\'s on your mind')
        cy.get('input').contains("Share")
        cy.get('input[name="text"]')

        // ---------------- CLEANUP ---------------- //
        
        cy.get('a').contains('sign out').click()
    })

    it('Check if message is in private timeline', () => {

        // ---------------- SETUP ---------------- //
        
        const message = uid()
        const username = uid()

        // register
        cy.register(username, "user@example.com", "password", "password")

        // login
        cy.login(username, "password")
        
        // create message
        cy.createMessage(username, message)

        // ---------------- TEST ---------------- //

        cy.contains('li', username).contains(message)

        // ---------------- CLEANUP ---------------- //

        cy.get('a').contains('sign out').click()
    })

    it('Check if message is in public timeline', () => {

        // ---------------- SETUP ---------------- //

        const message = uid()
        const username = uid()

        // register
        cy.register(username, "user@example.com", "password", "password")

        // login
        cy.login(username, "password")
        
        // create message
        cy.createMessage(username, message)

        cy.visit('localhost:8080/public')

        // ---------------- TEST ---------------- //

        cy.contains('li', username).contains(message)

        // ---------------- CLEANUP ---------------- //

        cy.get('a').contains('sign out').click()
    })

    it('Check if message is in users timeline', () => {

        // ---------------- SETUP ---------------- //

        const message = uid()
        const username = uid()

        // register
        cy.register(username, "user@example.com", "password", "password")

        // login
        cy.login(username, "password")
        
        // create message
        cy.createMessage(username, message)

        cy.get('a').contains('sign out').click()

        cy.visit('localhost:8080/'+username)

        // ---------------- TEST ---------------- //

        cy.contains('li', username).contains(message)

    })

    it('Check if messages are in correct order', () => {

        // ---------------- SETUP ---------------- //

        const message1 = uid()
        const message2 = uid()
        const username = uid()

        // register
        cy.register(username, "user@example.com", "password", "password")

        // login
        cy.login(username, "password")
        
        // create first message
        cy.createMessage(username, message1)

        // create second message
        cy.createMessage(username, message2)


        // ---------------- TEST ---------------- //

        // See if messages are in order
        cy.get('.messages > li').eq(0).should('contain', message2)
        cy.get('.messages > li').eq(1).should('contain', message1)

        // ---------------- CLEANUP ---------------- //

        cy.get('a').contains('sign out').click()

    })
})
