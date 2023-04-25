import { uid } from 'uid';

describe('Logout logs user out', () => {    
    it('Checks if logout logs the user out', () => {

        // ---------------- SETUP ---------------- //

        const username = uid()

        // register user
        cy.register(username, "user@example.com", "password", "password")


        // login
        cy.login(username, "password")

        
        // sign out
        cy.get('a').contains('sign out').click()

        // ---------------- TEST ---------------- //

        // make sure the buttons that appear are the ones that only are there for signed out users
        cy.get('a').contains('public timeline')
        cy.get('a').contains('sign up')
        cy.get('a').contains('sign in')
    })
})