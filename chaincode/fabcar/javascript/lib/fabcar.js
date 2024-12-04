/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Kalpcontract, Kalpsdk } = require('kalp-sdk-node');
// const { Kalpsdk } = require('./klap ');

class FabCar extends Kalpcontract {

    constructor() {
        console.info('============= START : FabCar constructor ===========');
        super('Myfabcar', true);
      }

       async createCar(ctx, carData) {
        console.info('============= START : Create Car ===========');

        let input = JSON.parse(carData)
        console.info('input',input);


        let carNumber  = input.CarNumber
        console.info('carNumber',carNumber);

        await ctx.putStateWithKYC(carNumber, Buffer.from(JSON.stringify(input)));
        console.info('============= END : Create Car ===========');
    }

     async createCarwithGasFee(ctx, carData) {
        console.info('============= START : Create Car ===========');

        let input = JSON.parse(carData)
        console.info('input',input);


        let carNumber  = input.CarNumber
        console.info('carNumber',carNumber);

        await ctx.putStateWithKYC(carNumber, Buffer.from(JSON.stringify(input)));
        console.info('============= END : Create Car ===========');
    }
    // async createCar(ctx, carNumber, make, model, color, owner) {
    //     console.info('============= START : Create Car ===========');

    //     const car = {
    //         color,
    //         docType: 'car',
    //         make,
    //         model,
    //         owner,
    //     };

    //     await Kalpsdk.putStateWithoutKYC(carNumber, Buffer.from(JSON.stringify(car)));
    //     console.info('============= END : Create Car ===========');
    // }

    async queryAllCars(ctx) {
        const startKey = '';
        const endKey = '';
        const allResults = [];
        for await (const {key, value} of ctx.stub.getStateByRange(startKey, endKey)) {
            const strValue = Buffer.from(value).toString('utf8');
            let record;
            try {
                record = JSON.parse(strValue);
            } catch (err) {
                console.log(err);
                record = strValue;
            }
            allResults.push({ Key: key, Record: record });
        }
        console.info(allResults);
        return JSON.stringify(allResults);
    }

    async changeCarOwner(ctx, carNumber, newOwner) {
        console.info('============= START : changeCarOwner ===========');

        const carAsBytes = await ctx.stub.getState(carNumber); // get the car from chaincode state
        if (!carAsBytes || carAsBytes.length === 0) {
            throw new Error(`${carNumber} does not exist`);
        }
        const car = JSON.parse(carAsBytes.toString());
        car.owner = newOwner;

        await ctx.stub.putState(carNumber, Buffer.from(JSON.stringify(car)));
        console.info('============= END : changeCarOwner ===========');
    }

}

module.exports = FabCar;
