# Fast & Scalable Decision Tree (C4.5)

## Overview
This project implements a high-performance and scalable Decision Tree (C4.5) classifier in Go. Decision Trees are a fundamental machine learning technique used for classification tasks. The solution is designed to efficiently handle large datasets with minimal memory overhead.

## Problem Statement
Your challenge is to implement a C4.5 decision tree classifier that is efficient, modular, and capable of handling large datasets with minimal memory overhead. Bonus points are awarded for implementing parallelization.

## Features
- **Training**: Construct a decision tree from labeled data.
- **Prediction**: Make predictions using the trained decision tree model.
- **Pruning**: Reduce overfitting by pruning the decision tree based on a validation set.
- **Error Handling**: Robust error handling for common issues such as missing files or incorrect input.

## Command-Line Interface (CLI) Documentation

### Commands
The `dt` command-line tool provides two primary commands:

1. **Train**: Constructs a decision tree from labeled data and saves the trained model.
2. **Predict**: Uses the trained model to make predictions on new data.

### 1. Training a Decision Tree
