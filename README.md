# imagerepo
a Shopify backend developer intern application task (Image Repository)

## Base url
[click here](https://shopify-img-repo.herokuapp.com)

## Prerequisite

- Postman

## Brief details to note

- You have the option to upload an image while you have an account on the image repository or not.

- If you do not have an account, you can only upload image with `public` and not `private` permission, and it will be visible to the whole world. you can also only search for `public` images

- If you have an account, you can upload both images with `public` and `private` permissions set to the image. you can also search for all `public` images and `private` images you own.

## How to use the app

- To register a account on the image repository, supply your email and password.

- To login to your account, supply your email and password.

- To add image, you have the option to add one or multiple images, using postman, you supply (image, text, image type (private or public and image characteristics)).

- To search for image, you have the option to search all images, search by text, search by images characteristics, search by another image in the repository using id of the image.

- More details of how application is used in Postman can be found in documentation below 😇

## Documentation
[find here](https://documenter.getpostman.com/view/4823089/TVzVibow)

## Endpoints

- POST Register /register

- POST Login /login

- GET Search Images /search

- POST Add Images /image

