import { uid } from 'uid';

describe('Follow', () => {
    
    it('User follows and message shows up', () => {

        // ---------------- SETUP ---------------- //

        const message = uid()
        const username1 = uid()
        const username2 = uid()

        // register username1
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username1)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // login as username1
        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username1)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        cy.visit('localhost:8080/'+username1)

        // create message as username1
        cy.get('input[name="text"]').type(message)
        cy.get('input').contains("Share").click()

        cy.get('a').contains('sign out').click()

        // register username2
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username2)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // login as username2
        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username2)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        // follow username1 as username2
        cy.visit('localhost:8080/'+username1)
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
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username1)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // login as username1
        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username1)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()

        // create message as username1
        cy.visit('localhost:8080/'+username1)

        cy.get('input[name="text"]').type(message)
        cy.get('input').contains("Share").click()

        cy.get('a').contains('sign out').click()

        // register username2
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username2)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // login as username2
        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username2)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        // follow username1 as username2
        cy.visit('localhost:8080/'+username1)
        cy.get('.follow').click()

        // unfollow username1 as username2
        cy.visit('localhost:8080/'+username1)
        cy.get('.unfollow').click()

        cy.visit('localhost:8080/')

        // ---------------- TEST ---------------- //

        // message of username1 shouldnt be on home page of username2
        cy.contains('li', username1).should('not.exist')

        // ---------------- CLEANUP ---------------- //
        
        cy.get('a').contains('sign out').click()
    })

})
