import { uid } from 'uid';

describe('Follow', () => {

    it('User follows and message shows up', () => {

        // ---------------- SETUP ---------------- //

        const message = uid()
        const username1 = uid()
        const username2 = uid()

        // register username1
        cy.register(username1, "user@example.com", "password", "password")

        // login as username1
        cy.login(username1, "password")

        // create message as username1
        cy.createMessage(username1, message)

        cy.get('a').contains('sign out').click()

        // register username2
        cy.register(username2, "user@example.com", "password", "password")

        // login as username2
        cy.login(username2, "password")

        // follow username1 as username2
        cy.visit('localhost:8080/' + username1)
        cy.get('.follow').click()

        cy.visit('localhost:8080/')

        // ---------------- TEST ---------------- //

        // message of username1 shows up on main page of username2
        cy.contains('li', username1).contains(message)

        // ---------------- CLEANUP ---------------- //

        cy.get('a').contains('sign out').click()
    })

    it('User unfollows and message disappears', () => {

        // ---------------- SETUP ---------------- //

        const message = uid()
        const username1 = uid()
        const username2 = uid()

        // register username1
        cy.register(username1, "user@example.com", "password", "password")

        // login as username1
        cy.login(username1, "password")

        // create message as username1
        cy.createMessage(username1, message)

        cy.get('a').contains('sign out').click()

        // register username2
        cy.register(username2, "user@example.com", "password", "password")

        // login as username2
        cy.login(username2, "password")


        // follow username1 as username2
        cy.visit('localhost:8080/' + username1)
        cy.get('.follow').click()

        // unfollow username1 as username2
        cy.visit('localhost:8080/' + username1)
        cy.get('.unfollow').click()

        cy.visit('localhost:8080/')

        // ---------------- TEST ---------------- //

        // message of username1 shouldnt be on home page of username2
        cy.contains('li', username1).should('not.exist')

        // ---------------- CLEANUP ---------------- //

        cy.get('a').contains('sign out').click()
    })

})
