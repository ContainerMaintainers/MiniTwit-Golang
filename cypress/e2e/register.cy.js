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

        // click on sign up button
        cy.get('input').contains('Sign Up').click()

        // ---------------- TEST ---------------- //

        // check if error message shows up
        cy.contains("Error: All fields are required")
    })

    it('Checks for "you have to enter a valid email address" error message', () => {

        // ---------------- SETUP ---------------- //

        // register with invalid email
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(uid())
        cy.get('input[name="email"]').type('example')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // ---------------- TEST ---------------- //

        // check if error message shows up       
        cy.contains("Error: You have to enter a valid email address")
    })

    it('Checks for "the two passwords do not matchs" error message', () => {
        
        // ---------------- SETUP ---------------- //

        // register with non-matching passwords
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(uid())
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('wrong')
        cy.get('input').contains('Sign Up').click()

        // ---------------- TEST ---------------- //

        // check if error message shows up
        cy.contains("Error: The two passwords do not match")
    })

    it('Checks for "the username is already taken" error message', () => {
        
        // ---------------- SETUP ---------------- //

        const username = uid()

        // register user
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // register the same user
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // ---------------- TEST ---------------- //

        // check if error message shows up
        cy.contains("Error: The username is already taken")
    })
})

describe('Redirection', () => {
    it('Checks that register redirects to login page', () => {

        // ---------------- SETUP ---------------- //

        // register
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(uid())
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

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
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // ---------------- TEST ---------------- //

        // check that the users own page exists
        cy.visit('localhost:8080/' + username)
    })
})