import { uid } from 'uid';

describe('Logout logs user out', () => {    
    it('Checks if logout logs the user out', () => {

        // ---------------- SETUP ---------------- //

        const username = uid()

        // register user
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // login
        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        // sign out
        cy.get('a').contains('sign out').click()

        // ---------------- TEST ---------------- //

        // make sure the buttons that appear are the ones that only are there for signed out users
        cy.get('a').contains('public timeline')
        cy.get('a').contains('sign up')
        cy.get('a').contains('sign in')
    })
})