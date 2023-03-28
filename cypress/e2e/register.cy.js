import { uid } from 'uid';

describe('Register page responds', () => {
    it('Checks if /register responds', () => {
        cy.visit('localhost:8080/register')
    })
})

describe('Register error messages', () => {
    it('Checks for "all fields are required" error message', () => {

        // ---------------- SETUP ---------------- //

        cy.visit('localhost:8080/register')
        cy.get('input').contains('Sign Up').click()

        // ---------------- TEST ---------------- //

        // check if error message shows up
        cy.contains("Error: All fields are required")
    })

    it('Checks for "you have to enter a valid email address" error message', () => {

        // ---------------- SETUP ---------------- //
        const username = uid()

        // register with invalid email
        cy.register(username, "example.com", "password", "password")

        // ---------------- TEST ---------------- //

        // check if error message shows up       
        cy.contains("Error: You have to enter a valid email address")
    })

    it('Checks for "the two passwords do not matchs" error message', () => {
        
        // ---------------- SETUP ---------------- //
        const username = uid()

        // register with non-matching passwords
        cy.register(username, "user@example.com", "password1", "password2")

        // ---------------- TEST ---------------- //

        // check if error message shows up
        cy.contains("Error: The two passwords do not match")
    })

    it('Checks for "the username is already taken" error message', () => {
        
        // ---------------- SETUP ---------------- //

        const username = uid()

        // register user
        cy.register(username, "user@example.com", "password", "password")

        // register the same user
        cy.register(username, "user@example.com", "password", "password")

        // ---------------- TEST ---------------- //

        // check if error message shows up
        cy.contains("Error: The username is already taken")
    })
})

describe('Redirection', () => {
    it('Checks that register redirects to login page', () => {

        // ---------------- SETUP ---------------- //
        const username = uid()

        // register
        cy.register(username, "user@example.com", "password", "password")

        // ---------------- TEST ---------------- //

        // check that the current page is the /login page
        cy.url().should('contain', '/login')
    })
})

describe('User exists', () => {
    it('Checks that registered user exists', () => {

        // ---------------- SETUP ---------------- //

        const username = uid()

        // register
        cy.register(username, "user@example.com", "password", "password")

        // ---------------- TEST ---------------- //

        // check that the users own page exists
        cy.visit('localhost:8080/' + username)
    })
})