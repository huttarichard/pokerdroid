/**
 * Copyright 1993-2012 NVIDIA Corporation.  All rights reserved.
 *
 * Please refer to the NVIDIA end user license agreement (EULA) associated
 * with this source code for terms and conditions that govern your use of
 * this software. Any use, reproduction, disclosure, or distribution of
 * this software and related documentation outside the terms of the EULA
 * is strictly prohibited.
 */
#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <chrono>
#include <thread>
#include "HandEvaluator.h"
#include "EquityCalculator.h"

using namespace std;
using namespace omp;
using namespace std::chrono;


//may return 0 when not able to detect
const auto processor_count = std::thread::hardware_concurrency();

extern "C" double hand_equity(const char *hole_card, const char *community_card, const char *dead_card, const int nb_players,
		const int nb_board_cards_ = 5, const double std_tol = 1e-5)
{

	high_resolution_clock::time_point t1 = high_resolution_clock::now();
	unsigned nb_board_cards = nb_board_cards_;
	EquityCalculator eq;
	vector<CardRange> ranges (nb_players);
	ranges[0] = hole_card;
	for( int a = 1; a < nb_players; a = a + 1 ) {
		ranges[a] = "random";
	}


	uint64_t board = CardRange::getCardMask(community_card); //
	uint64_t dead = CardRange::getCardMask(dead_card); //
	// stop when standard error below std_tol%

	auto callback = [&eq](const EquityCalculator::Results& results) {
			//cout << results.equity[0] << " " << 100 * results.progress
			//	<< " " << 1e-6 * results.intervalSpeed << endl;

		if (results.time > 1) // Stop after 1s
			eq.stop();
	};
	double updateInterval = 0.0001; // Callback called every 0.0001s.
	unsigned threads = processor_count; // max hardware parallelism (default)
	eq.start(ranges, board, dead, false, std_tol, callback, updateInterval, threads, nb_board_cards);
	eq.wait();
	auto r = eq.getResults();

    high_resolution_clock::time_point t2 = high_resolution_clock::now();
	auto duration = duration_cast<microseconds>( t2 - t1 ).count();

	// cout << "[From C++] Hole cards:" << hole_card <<endl;
	// cout << "[From C++] Community cards:" << community_card <<endl;

    // cout << "[From C++] Total time taken: "<< 1e-3*duration << " [ms], Time taken for evals: " << 1e3*r.time << " [ms]"<<endl;
    // cout << "[From C++] Speed: " << 1e-6 * r.speed << " [Meval/s], std: " << r.stdev << endl;

	// cout << "[From C++] Cores: " << processor_count <<endl;
	// cout << "[From C++] Wins: " <<  r.wins[0] <<endl;
	// cout << "[From C++] Ties: " <<  r.ties[0] <<endl;
	// cout << "[From C++] Evaulations: " <<  r.evaluations <<endl;
	// cout << "[From C++] Equity: " <<  r.equity[0] <<endl;
	return r.equity[0];
}



extern "C" void hand_equity_with_draw(const char *hole_card, const char *community_card, const char *dead_card, 
                                      const int nb_players, double *win_equity, double *draw_equity,
                                      const int nb_board_cards_ = 5, const double std_tol = 1e-5)
{
    high_resolution_clock::time_point t1 = high_resolution_clock::now();
    unsigned nb_board_cards = nb_board_cards_;
    EquityCalculator eq;
    vector<CardRange> ranges (nb_players);
    ranges[0] = hole_card;
    for( int a = 1; a < nb_players; a = a + 1 ) {
        ranges[a] = "random";
    }

    uint64_t board = CardRange::getCardMask(community_card);
    uint64_t dead = CardRange::getCardMask(dead_card);

    auto callback = [&eq](const EquityCalculator::Results& results) {
        if (results.time > 1) // Stop after 1s
            eq.stop();
    };
    
    double updateInterval = 0.0001;
    unsigned threads = processor_count;
    eq.start(ranges, board, dead, false, std_tol, callback, updateInterval, threads, nb_board_cards);
    eq.wait();
    auto r = eq.getResults();

    *win_equity = r.wins[0] / (double)r.hands;
    *draw_equity = r.ties[0] / (double)r.hands;
}